package parser

import (
	"testing"
)

func TestAlbLogParser(t *testing.T) {
	cases := []struct {
		line, method, datetime, protcol, path string
	}{
		{
			"http 2017-09-09T00:55:00.228196Z app/albname/12a34bc6d78e9f0 111.222.10.240:57965 192.168.131.118:80 0.020 0.019 0.001 200 200 623 111603 \"GET http://wwww.host.com:80/path/?a=1&b=2&c=3 HTTP/1.1\" \"Mozilla/5.0 (iPhone; CPU iPhone OS 10_2_1 like Mac OS X) AppleWebKit/602.4.6 (KHTML, like Gecko) Version/10.0 Mobile/14D27 Safari/602.1\" - - arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/albname/abcdefgh12345678 \"Root=1-11111be1-1110300f111111111aa1a1aa\"",
			"GET",
			"2017-09-09T00:55:00.228196Z",
			"http",
			"path/?a=1&b=2&c=3",
		},
	}

	p := NewALBLogParser()
	for _, c := range cases {
		entry, err := p.Parse(c.line)
		if err != nil {
			t.Fatalf("parse error:%#v\n", err)
		}
		if entry.Method != c.method {
			t.Fatalf("parse error Method:%s\n", entry.Method)
		}
		if entry.Protocol != c.protcol {
			t.Fatalf("parse error Protocol:%s\n", entry.Protocol)
		}
		if entry.DateTime.Format(p.format) != c.datetime {
			t.Fatalf("parse error Time:%v\n", entry.DateTime)
		}
		if entry.Path != c.path {
			t.Fatalf("parse error Path:%s\n", entry.Path)
		}
	}
}
