package jsonschema

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/akamensky/base58"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

func isString(s ...interface{}) (interface{}, error) {
	var params ParamMap
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if len(s) > 1 {
		params = s[1].(ParamMap)
	}
	if q, ok := s[0].(string); ok {
		if err := checkStrConstraints(q, params); err != nil {
			return nil, err
		}
		return q, nil
	}
	return nil, ErrViolationType
}

func isFloat(s ...interface{}) (interface{}, error) {
	var params ParamMap
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if len(s) > 1 {
		params = s[1].(ParamMap)
	}
	var q float64
	switch n := s[0].(type) {
	case float32:
		q = float64(n)
	case float64:
		q = n
	case int:
		q = float64(n)
	default:
		return nil, ErrViolationType
	}
	if min, ok, err := params.AsFloat("min"); err != nil {
		return nil, err
	} else if ok {
		if q < min {
			return nil, ErrParamConstraint
		}
	}
	if max, ok, err := params.AsFloat("max"); err != nil {
		return nil, err
	} else if ok {
		if q > max {
			return nil, ErrParamConstraint
		}
	}
	return q, nil

}

func isInt(s ...interface{}) (interface{}, error) {
	var params ParamMap
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if len(s) > 1 {
		params = s[1].(ParamMap)
	}
	var q int
	switch n := s[0].(type) {
	case float32:
		if s[0] == float32(int(n)) {
			q = int(n)
		} else {
			return nil, ErrViolationType
		}
	case float64:
		if s[0] == float64(int(n)) {
			q = int(n)
		} else {
			return nil, ErrViolationType
		}
	case int:
		q = n
	default:
		return nil, ErrViolationType
	}

	if min, ok, err := params.AsInt("min"); err != nil {
		return nil, err
	} else if ok {
		if int64(q) < min {
			return nil, ErrParamConstraint
		}
	}
	if max, ok, err := params.AsInt("max"); err != nil {
		return nil, err
	} else if ok {
		if int64(q) > max {
			return nil, ErrParamConstraint
		}
	}
	return q, nil
}

func isDir(s ...interface{}) (interface{}, error) {
	if len(s) != 1 {
		return nil, ErrViolationType
	}
	if str, ok := s[0].(string); ok {
		p := path.Clean(str)
		if stat, err := os.Stat(p); err != nil {
			return nil, err
		} else if !stat.IsDir() {
			return nil, ErrViolationType
		}
		return p, nil
	}
	return nil, ErrViolationType
}

func isFile(s ...interface{}) (interface{}, error) {
	if len(s) != 1 {
		return nil, ErrViolationType
	}
	if str, ok := s[0].(string); ok {
		p := path.Clean(str)
		if stat, err := os.Stat(p); err != nil {
			return nil, err
		} else if stat.IsDir() {
			return nil, ErrViolationType
		}
		return p, nil
	}
	return nil, ErrViolationType
}

func isDuration(s ...interface{}) (interface{}, error) {
	var err error
	var params ParamMap
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if len(s) > 1 {
		params = s[1].(ParamMap)
	}
	var q time.Duration
	if d, ok := s[0].(string); !ok {
		return nil, ErrViolationType
	} else if q, err = time.ParseDuration(d); err != nil {
		return nil, ErrViolationType
	}
	if minS, ok, err := params.AsString("min"); err != nil {
		return nil, err
	} else if ok {
		if min, err := time.ParseDuration(minS); err != nil {
			return nil, err
		} else if q < min {
			return nil, ErrParamConstraint
		}
	}
	if maxS, ok, err := params.AsString("max"); err != nil {
		return nil, err
	} else if ok {
		if max, err := time.ParseDuration(maxS); err != nil {
			return nil, err
		} else if q > max {
			return nil, ErrParamConstraint
		}
	}
	return q, nil
}

func checkStrLen(s string, params ParamMap) error {
	if l, ok, err := params.AsInt("len"); err != nil {
		return err
	} else if ok {
		if len(s) != int(l) {
			return ErrParamConstraint
		}
	}
	return nil
}

func checkStrMinLen(s string, params ParamMap) error {
	if l, ok, err := params.AsInt("min"); err != nil {
		return err
	} else if ok {
		if len(s) < int(l) {
			return ErrParamConstraint
		}
	}
	return nil
}

func checkStrMaxLen(s string, params ParamMap) error {
	if l, ok, err := params.AsInt("max"); err != nil {
		return err
	} else if ok {
		if len(s) > int(l) {
			return ErrParamConstraint
		}
	}
	return nil
}

func checkStrConstraints(s string, params ParamMap) error {
	if err := checkStrLen(s, params); err != nil {
		return err
	}
	if err := checkStrMinLen(s, params); err != nil {
		return err
	}
	if err := checkStrMaxLen(s, params); err != nil {
		return err
	}
	return nil
}

func isEncoded(removePrefix bool, decodeFunc func(string) error, s ...interface{}) (interface{}, error) {
	var params ParamMap
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if len(s) > 1 {
		params = s[1].(ParamMap)
	}
	if str, ok := s[0].(string); ok {
		if removePrefix {
			if strings.HasPrefix(str, "x") {
				str = str[1:]
			} else if strings.HasPrefix(str, "0x") {
				str = str[2:]
			}
		}
		if err := decodeFunc(str); err != nil {
			return nil, err
		}
		if err := checkStrConstraints(str, params); err != nil {
			return nil, err
		}
		return str, nil
	}
	return nil, ErrViolationType
}

func isHex(s ...interface{}) (interface{}, error) {
	return isEncoded(true, func(s string) error {
		_, err := hex.DecodeString(s)
		return err
	}, s...)
}

func isBase64(s ...interface{}) (interface{}, error) {
	return isEncoded(false, func(s string) error {
		_, err := base64.StdEncoding.DecodeString(s)
		return err
	}, s...)
}

func isBase58(s ...interface{}) (interface{}, error) {
	return isEncoded(false, func(s string) error {
		_, err := base58.Decode(s)
		return err
	}, s...)
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

func isIPv4Addr(s ...interface{}) (interface{}, error) {
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if str, ok := s[0].(string); ok {
		ip := net.ParseIP(str)
		if !isIPv4(ip) {
			return nil, ErrViolationType
		}
		return ip.String(), nil
	}
	return nil, ErrViolationType
}

func isIPv6Addr(s ...interface{}) (interface{}, error) {
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if str, ok := s[0].(string); ok {
		ip := net.ParseIP(str)
		if !isIPv6(ip) {
			return nil, ErrViolationType
		}
		return ip.String(), nil
	}
	return nil, ErrViolationType
}

func isIPv4Net(s ...interface{}) (interface{}, error) {
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if str, ok := s[0].(string); ok {
		ip, ipnet, err := net.ParseCIDR(str)
		if err != nil {
			return nil, err
		}
		if ipnet == nil {
			return nil, err
		}
		if !isIPv4(ip) {
			return nil, ErrViolationType
		}
		return ip.String(), nil
	}
	return nil, ErrViolationType
}

func isIPv6Net(s ...interface{}) (interface{}, error) {
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if str, ok := s[0].(string); ok {
		ip, ipnet, err := net.ParseCIDR(str)
		if err != nil {
			return nil, err
		}
		if ipnet == nil {
			return nil, err
		}
		if !isIPv6(ip) {
			return nil, ErrViolationType
		}
		return ip.String(), nil
	}
	return nil, ErrViolationType
}

func isHostname(s ...interface{}) (interface{}, error) {
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if str, ok := s[0].(string); ok {
		if hn, err := os.Hostname(); err != nil {
			return nil, err
		} else if hn != str {
			return nil, ErrViolationType
		}
		return str, nil
	}
	return nil, ErrViolationType
}

func isNIC(s ...interface{}) (interface{}, error) {
	return isNICver(0, s...)
}

func isNIC4(s ...interface{}) (interface{}, error) {
	return isNICver(4, s...)
}

func isNIC6(s ...interface{}) (interface{}, error) {
	return isNICver(6, s...)
}

func isNICver(ver int, s ...interface{}) (interface{}, error) {
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if str, ok := s[0].(string); ok {
		inf, err := net.Interfaces()
		if err != nil {
			return nil, err
		}
		for _, i := range inf {
			if i.Name == str {
				if ver != 0 {
					addr, _ := i.Addrs()
					for _, a := range addr {
						ip, _, _ := net.ParseCIDR(a.String())
						if ver == 4 && isIPv4(ip) {
							return str, nil
						}
						if ver == 6 && isIPv6(ip) {
							return str, nil
						}
					}
					return nil, ErrViolationType
				}
				return str, nil
			}
		}
		return nil, ErrViolationType
	}
	return nil, ErrViolationType
}

func lookupAddr(ver int, s ...interface{}) (interface{}, error) {
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if str, ok := s[0].(string); ok {
		addrs, err := net.LookupHost(str)
		if err != nil {
			return nil, ErrViolationType
		}
		for _, addr := range addrs {
			ip := net.ParseIP(addr)
			if ver == 4 && isIPv4(ip) {
				return str, nil
			}
			if ver == 6 && isIPv6(ip) {
				return str, nil
			}
		}
	}
	return nil, ErrViolationType
}

func lookupIPv4(s ...interface{}) (interface{}, error) {
	return lookupAddr(4, s...)
}

func lookupIPv6(s ...interface{}) (interface{}, error) {
	return lookupAddr(6, s...)
}
