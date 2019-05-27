package libgobuster

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	//"github.com/google/uuid"
)

// GobusterDNS is the main type to implement the interface
type GobusterDNS struct{}

// Setup is the setup implementation of gobusterdns
func (d GobusterDNS) Setup(g *Gobuster) error {

	/*
		// Resolve a subdomain sthat probably shouldn't exist
		guid := uuid.New()
		wildcardIps, err := g.DNSLookup(fmt.Sprintf("%s.%s", guid, g.Opts.URL))
		if err == nil {
			g.IsWildcard = true
			g.WildcardIps.AddRange(wildcardIps)
			log.Printf("[-] Wildcard DNS found. IP address(es): %s", g.WildcardIps.Stringify())
			if !g.Opts.WildcardForced {
				return fmt.Errorf("To force processing of Wildcard DNS, specify the '-fw' switch.")
			}
		}
	*/

	if !g.Opts.Quiet {
		// Provide a warning if the base domain doesn't resolve (in case of typo)
		_, err := g.DNSLookup(g.Opts.URL)
		if err != nil {
			// Not an error, just a warning. Eg. `yp.to` doesn't resolve, but `cr.py.to` does!
			log.Printf("[-] Unable to validate base domain: %s", g.Opts.URL)
		}
	}

	return nil
}

// Process is the process implementation of gobusterdns
func (d GobusterDNS) Process(g *Gobuster, word string) ([]Result, error) {
	subdomain := fmt.Sprintf("%s.%s", word, g.Opts.URL)
	fmt.Println("try ", subdomain)
	var ret []Result

	ips, err := g.DNSLookup(subdomain) //return ip address contain ipv4 and ipv6
	if err != nil {
		if g.Opts.Verbose {
			fmt.Printf("lookup host %s error:%s\n", subdomain, err)
		}
		return nil, err
	}
	if g.Opts.ShowIPs {
		for _, ip := range ips {
			if ip.Contains(".") { //ipv4
				if g.Opts.ShowA {
					result := Result{
						Entity:  subdomain,
						Extra:   ip,
						DnsType: "A",
					}
					ret = append(ret, result)
				}
			} else { //ipv6
				if g.Opts.ShowAAAA {
					result := Result{
						Entity:  subdomain,
						Extra:   ip,
						DnsType: "AAAA",
					}
					ret = append(ret, result)
				}
			}
		}
	}
	if g.Opts.ShowCNAME {
		cname, err := g.DNSLookupCname(subdomain)
		if err == nil {
			result := Result{
				Entity:  subdomain,
				Extra:   cname,
				DnsType: "CNAME",
			}
			ret = append(ret, result)
		} else {
			fmt.Printf("lookup CNAME %s error:%s\n", subdomain, err)
		}
	}
	if g.Opts.ShowMX {
		mxs, err := g.DNSLookupMX(subdomain)
		if err == nil {
			for _, mx = range mxs {
				result := Result{
					Entity:  subdomain,
					Extra:   mx,
					DnsType: "MX",
				}
			}
		} else {
			fmt.Printf("lookup MX %s error:%s\n", subdomain, err)
		}
	}

	return ret, nil
}

// ResultToString is the to string implementation of gobusterdns
func (d GobusterDNS) ResultToString(g *Gobuster, r *Result) (*string, error) {
	buf := &bytes.Buffer{}

	if r.Status == 404 {
		if _, err := fmt.Fprintf(buf, "Missing: %s\n", r.Entity); err != nil {
			return nil, err
		}
	} else {
		if _, err := fmt.Fprintf(buf, "Found: %s [%s] [%s]\n", r.Entity, r.Extra, r.DnsType); err != nil {
			return nil, err
		}
	}

	s := buf.String()
	return &s, nil
}
