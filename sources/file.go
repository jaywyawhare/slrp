package sources

import (
    "bufio"
    "context"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/nfx/slrp/pmux"
)

// ip:port or ip:port protocol (protocol one of http, https, socks4, socks5) along with username and password
func FileSource(filePath string) Source {
    clean := filepath.Clean(filePath)
    return Source{
        ID:        1000,
        name:      "file",
        Homepage:  clean,
        Frequency: 0,
        Seed:      true,
        Feed: simpleGen(func(context.Context, *http.Client) ([]pmux.Proxy, error) {
            f, err := os.Open(clean)
            if err != nil {
                return nil, err
            }
            defer f.Close()
            var out []pmux.Proxy
            scanner := bufio.NewScanner(f)
            for scanner.Scan() {
                line := strings.TrimSpace(scanner.Text())
                if line == "" || strings.HasPrefix(line, "#") {
                    continue
                }
                proto := "http"
                addr := line
                var user, pass string
                fields := strings.Fields(line)
                switch len(fields) {
                case 0:
                    continue
                case 1:
                    parts := strings.Split(fields[0], ":")
                    if len(parts) >= 4 {
                        addr = strings.Join(parts[0:2], ":")
                        user = parts[2]
                        pass = parts[3]
                    } else {
                        addr = fields[0]
                    }
                default:
                    proto = strings.ToLower(fields[len(fields)-1])
                    left := strings.Join(fields[:len(fields)-1], " ")
                    parts := strings.Split(left, ":")
                    if len(parts) >= 4 {
                        addr = strings.Join(parts[0:2], ":")
                        user = parts[2]
                        pass = parts[3]
                    } else {
                        addr = left
                    }
                }
                allowed := map[string]struct{}{"http": {}, "https": {}, "socks4": {}, "socks5": {}}
                if _, ok := allowed[proto]; !ok {
                    continue
                }
                p := pmux.NewProxy(addr, proto)
                if p != 0 {
                    if user != "" || pass != "" {
                        pmux.SetProxyAuth(p, user, pass)
                    }
                    out = append(out, p)
                }
            }
            if err := scanner.Err(); err != nil {
                return nil, err
            }
            return out, nil
        }),
    }
}


