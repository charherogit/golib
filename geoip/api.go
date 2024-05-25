package geoip

import (
	"net"
)

// 如: CN/US/JP
func (i *IpSearch) GetCountryIsoCode(val string) (string, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	ip := net.ParseIP(val)
	record, err := i.reader.City(ip)
	if err != nil {
		return "", err
	}
	return record.Country.IsoCode, nil
}

// 如: 中国/美国/日本
func (i *IpSearch) GetCountryName(val string) (string, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	ip := net.ParseIP(val)
	record, err := i.reader.City(ip)
	if err != nil {
		return "", err
	}
	return record.Country.Names["zh-CN"], nil
}

// 如: 北京/纽约/东京
func (i *IpSearch) GetCityName(val string) (string, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	ip := net.ParseIP(val)
	record, err := i.reader.City(ip)
	if err != nil {
		return "", err
	}
	return record.City.Names["zh-CN"], nil
}
