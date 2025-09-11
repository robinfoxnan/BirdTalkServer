package core

import (
	"fmt"
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

// go get github.com/oschwald/geoip2-golang/v2
// GeoIPHelper 封装 GeoLite2-City 数据库功能
type GeoIPHelper struct {
	db   *geoip2.Reader
	lock sync.Mutex
}

// NewGeoIPHelper 创建 GeoIPHelper 实例
func NewGeoIPHelper(dbPath string) (*GeoIPHelper, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %v", err)
	}

	return &GeoIPHelper{
		db: db,
	}, nil
}

// Close 关闭数据库
func (g *GeoIPHelper) Close() {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.db != nil {
		g.db.Close()
		g.db = nil
	}
}

// CityInfo 包含 IP 对应的地址信息
type CityInfo struct {
	IP        string
	Country   string
	Province  string
	City      string
	Latitude  float64
	Longitude float64
}

// GetCityByIP 根据 IP 获取城市信息
func (g *GeoIPHelper) GetCityByIP(ipStr string) (*CityInfo, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.db == nil {
		return nil, fmt.Errorf("数据库未加载")
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("无效 IP: %s", ipStr)
	}

	record, err := g.db.City(ip)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	country := ""
	if name, ok := record.Country.Names["zh-CN"]; ok {
		country = name
	}

	province := ""
	if len(record.Subdivisions) > 0 {
		if name, ok := record.Subdivisions[0].Names["zh-CN"]; ok {
			province = name
		}
	}

	city := ""
	if name, ok := record.City.Names["zh-CN"]; ok {
		city = name
	}

	return &CityInfo{
		IP:        ipStr,
		Country:   country,
		Province:  province,
		City:      city,
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
	}, nil
}
