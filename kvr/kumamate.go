package kvr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	cacheFilePathKey = "KVR_CACHE_FILE"
	url              = "https://kumamate.net/vip/"
	rankSelector     = "#ctn-main > div.sbuu.data > div > table > tbody > tr > td"
)

type Client struct {
	client    *http.Client
	cachePath string
	useCache  bool
}

func New(hc *http.Client) (*Client, error) {
	if hc == nil {
		return nil, fmt.Errorf("http.Client is nil")
	}

	cachePath := os.Getenv(cacheFilePathKey)
	if len(cachePath) == 0 {
		return nil, fmt.Errorf("Environment variable '%s' is empty.", cacheFilePathKey)
	}

	return &Client{
		client:    hc,
		cachePath: cachePath,
		useCache:  true,
	}, nil
}

func (c *Client) NoCache() *Client {
	c.useCache = false
	return c
}

func (c *Client) ListVIPRank() (VIPRankList, error) {
	var vrs VIPRankList
	var err error
	if c.useCache {
		vrs, err = c.getVIPRankFromCache()
		if err == nil {
			return vrs, nil
		}

		// if err is not nil, fetch from origin
		fmt.Printf("Failed to get from cache. Try fetch. err:%v", err)
	}

	if !c.useCache || err != nil {
		vrs, err := c.fetchVIPRankFromOrigin()
		if err != nil {
			return nil, err
		}

		// save to cache
		bytes, err := json.Marshal(vrs)
		if err != nil {
			fmt.Printf("Failed to save to cache. err:%v", err)
			return vrs, nil
		}

		f, err := os.Create(c.cachePath)
		if err != nil {
			fmt.Printf("Failed to save to cache. err:%v", err)
			return vrs, nil
		}

		defer f.Close()

		if _, err := f.Write(bytes); err != nil {
			fmt.Printf("Failed to save to cache. err:%v", err)
			return vrs, nil
		}

		return vrs, nil
	}

	return c.getVIPRankFromCache()
}

func (c *Client) getVIPRankFromCache() (VIPRankList, error) {
	f, err := os.Open(c.cachePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bytes := bytes.Buffer{}
	if _, err = bytes.ReadFrom(f); err != nil {
		return nil, err
	}

	vrs := VIPRankList{}
	err = json.Unmarshal(bytes.Bytes(), &vrs)
	if err != nil {
		return nil, err
	}

	if len(vrs) == 0 {
		return nil, errors.New("Cache file is empty")
	}

	return vrs, nil
}

func (c *Client) fetchVIPRankFromOrigin() (VIPRankList, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch original data. statusCode:%d", res.StatusCode)
	}

	content, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var vrs VIPRankList = VIPRankList{}
	var current VIPRank
	content.Find(rankSelector).Each(func(i int, s *goquery.Selection) {
		if i < 3 {
			return
		}

		text := s.Text()

		switch i % 3 {
		case 0: // Rank ID

		case 1: // Title
			current.Title = text
		case 2: // Power Value
			pow, err := parsePower(text)
			if err != nil {
				fmt.Printf("Failed to parse power value. pow:%s err:%v", text, err)
				return
			}

			current.Min = pow
			tmp := current
			fmt.Printf("%#+v\n", tmp)
			vrs = append(vrs, tmp)
			current.Max = pow - 1
		}
	})

	return vrs, nil
}

var numReg = regexp.MustCompile("[^0-9]")

func parsePower(powStr string) (uint64, error) {
	powStr = numReg.ReplaceAllString(powStr, "")

	return strconv.ParseUint(powStr, 0, 64)
}
