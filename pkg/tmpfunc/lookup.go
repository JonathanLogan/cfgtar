package tmpfunc

import "net"

func ipLookup(ver int, s string) ([]string, error) {
	addr, err := net.LookupIP(s)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(addr))
	for _, a := range addr {
		if ver == 4 && isIPv4(a) {
			ret = append(ret, a.String())
		}
		if ver == 6 && isIPv6(a) {
			ret = append(ret, a.String())
		}
	}
	return ret, nil
}

func ipv4lookup(s string) ([]string, error) {
	return ipLookup(4, s)
}

func ipv6lookup(s string) ([]string, error) {
	return ipLookup(6, s)
}

func dnsTXT(s string) ([]string, error) {
	return net.LookupTXT(s)
}
