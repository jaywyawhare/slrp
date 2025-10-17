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

// ip:port or ip:port protocol (protocol one of http, https, socks4, socks5)
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
                if sp := strings.Fields(line); len(sp) >= 2 {
                    addr = sp[0]
                    proto = strings.ToLower(sp[1])
                }
                allowed := map[string]struct{}{"http": {}, "https": {}, "socks4": {}, "socks5": {}}
                if _, ok := allowed[proto]; !ok {
                    continue
                }
                p := pmux.NewProxy(addr, proto)
                if p != 0 {
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


