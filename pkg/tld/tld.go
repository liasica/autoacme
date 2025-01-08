// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-08, by liasica

package tld

import (
	"errors"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type List struct {
	prefix string // 主机记录
	base   string // 域名
	tld    string // 顶级域名
}

func Parse(str string) (l *List, err error) {
	str = strings.TrimSuffix(str, ".")
	var dn string
	dn, err = publicsuffix.EffectiveTLDPlusOne(str)
	if err != nil {
		return
	}
	i := strings.Index(dn, ".")
	if i < 0 {
		return nil, errors.New("tld: invalid domain")
	}

	tld := dn[i+1:]
	base := dn[:i]

	l = &List{
		base: base,
		tld:  tld,
	}
	if str != dn {
		l.prefix = str[:len(str)-len(dn)-1]
	}
	return
}

func (l *List) Domain() string {
	return l.base + "." + l.tld
}

// SubDomain 返回子域名 <AliyunDNS：主机记录>
func (l *List) SubDomain() string {
	return l.prefix
}
