package tmpfunc

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"text/template"
	"time"
)

var FuncMap = template.FuncMap{
	"hostname":    os.Hostname,
	"ipv4CIDR":    ipv4CIDR,
	"ipv4Mask":    ipv4Mask,
	"ipv6CIDR":    ipv6CIDR,
	"ipv6Mask":    ipv6Mask,
	"file":        file,
	"durationAs":  durationAs,
	"ipv4NICAddr": ipv4NICAddr,
	"ipv6NICAddr": ipv6NICAddr,
	"ipv4addr":    ipv4addr,
	"ipv6addr":    ipv6addr,
	"ipv4lookup":  ipv4lookup,
	"ipv6lookup":  ipv6lookup,
	"dnsTXT":      dnsTXT,
}

func asIP(ver int, s string) (string, error) {
	ip, _, err := net.ParseCIDR(s)
	if err != nil {
		return "", err
	}
	if ip == nil {
		return "", errors.New("not an IP address")
	}
	if ver == 4 && !isIPv4(ip) {
		return "", errors.New("not ipv4")
	}
	if ver == 6 && !isIPv6(ip) {
		return "", errors.New("not ipv4")
	}
	return ip.String(), nil
}

func ipv4addr(s string) (string, error) {
	return asIP(4, s)
}

func ipv6addr(s string) (string, error) {
	return asIP(6, s)
}

func isIPv4(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.To4() == nil {
		return false
	}
	return true
}

func isIPv6(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.To4() != nil {
		return false
	}
	if ip.To16() == nil {
		return false
	}
	return true
}

func ipNet(ver int, s string) (net.IPMask, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	if ver == 4 && !isIPv4(ip) {
		return nil, errors.New("not ipv4")
	}
	if ver == 6 && !isIPv6(ip) {
		return nil, errors.New("not ipv4")
	}
	if ipnet == nil {
		return nil, errors.New("no net")
	}
	return ipnet.Mask, nil
}

func ipv4CIDR(s string) (string, error) {
	if mask, err := ipNet(4, s); err != nil {
		return "", err
	} else {
		bit, _ := mask.Size()
		return strconv.FormatInt(int64(bit), 10), nil
	}
}

func ipv4Mask(s string) (string, error) {
	if mask, err := ipNet(4, s); err != nil {
		return "", err
	} else {
		return net.IP(mask).String(), nil
	}
}

func ipv6CIDR(s string) (string, error) {
	if mask, err := ipNet(6, s); err != nil {
		return "", err
	} else {
		bit, _ := mask.Size()
		return strconv.FormatInt(int64(bit), 10), nil
	}
}

func ipv6Mask(s string) (string, error) {
	if mask, err := ipNet(6, s); err != nil {
		return "", err
	} else {
		return net.IP(mask).String(), nil
	}
}

func file(s string) (string, error) {
	d, err := ioutil.ReadFile(s)
	return string(d), err
}

func durationAs(s ...string) (string, error) {
	dur := s[0]
	divider := time.Second
	if len(s) > 1 {
		dur = s[1]
		switch s[0] {
		case "s", "sec", "second":
			divider = time.Second
		case "m", "min", "minute":
			divider = time.Minute
		case "h", "hour":
			divider = time.Hour
		default:
			return "", errors.New("not hour/minute/second")
		}
	}
	durD, err := time.ParseDuration(dur)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(int64(durD/divider), 10), nil
}

func ipv4NICAddr(nic string) ([]string, error) {
	return nicNet(4, nic)
}

func ipv6NICAddr(nic string) ([]string, error) {
	return nicNet(6, nic)
}

func nicNet(ver int, nic string) ([]string, error) {
	inf, err := net.InterfaceByName(nic)
	if err != nil {
		return nil, err
	}
	addrs, err := inf.Addrs()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(addrs))
	for _, ad := range addrs {
		ip, _, err := net.ParseCIDR(ad.String())
		if err != nil {
			continue
		}
		if ver == 4 && isIPv4(ip) {
			ret = append(ret, ad.String())
		}
		if ver == 6 && isIPv6(ip) {
			ret = append(ret, ad.String())
		}
	}
	return ret, nil
}
