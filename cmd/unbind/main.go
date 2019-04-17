package main

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"regexp"
	"strconv"
)


func domainsToAddresses (domain string) (string, bool) {
	ok := false
	ip := ""
	re := regexp.MustCompile(`(ec2|ip)-(\d{1,3})-(\d{1,3})-(\d{1,3})-(\d{1,3})\.(compute|ec2).*`)
	if re.MatchString(domain) {
		ip = re.ReplaceAllString("ec2-54-173-113-15.compute-1.amazonaws.com", "$2.$3.$4.$5")
		ok = true
	}
	return ip, ok
}

type handler struct{}
func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := domainsToAddresses(domain)
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{ Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60 },
				A: net.ParseIP(address),
			})
		}
	}
	w.WriteMsg(&msg)
}

func main() {
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}