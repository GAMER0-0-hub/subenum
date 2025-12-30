package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

/* ================= COLORS ================= */

const (
    BLUE   = "\033[1;34m"
    GREEN  = "\033[1;32m"
    YELLOW = "\033[1;33m"
    CYAN   = "\033[1;36m"
    RED    = "\033[1;31m"
    RESET  = "\033[0m"
)

/* ================= FLAGS ================= */

var (
    domain string
    list   string
    outDir string
)

/* ================= SPINNER ================= */

func spinner(msg string, stop chan bool, wg *sync.WaitGroup) {
    defer wg.Done()
    frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
    i := 0
    for {
        select {
        case <-stop:
            fmt.Print("\r" + strings.Repeat(" ", len(msg)+5) + "\r")
            return
        default:
            fmt.Printf("\r%s %s", frames[i%len(frames)], msg)
            time.Sleep(90 * time.Millisecond)
            i++
        }
    }
}

/* ================= BANNER ================= */

func banner() {
    fmt.Println(CYAN)
    fmt.Println(`
        ███████╗██╗   ██╗██████╗ ███████╗███╗   ██╗██╗   ██╗███╗   ███╗
        ██╔════╝██║   ██║██╔══██╗██╔════╝████╗  ██║██║   ██║████╗ ████║
        ███████╗██║   ██║██████╔╝█████╗  ██╔██╗ ██║██║   ██║██╔████╔██║
        ╚════██║██║   ██║██╔══██╗██╔══╝  ██║╚██╗██║██║   ██║██║╚██╔╝██║
        ███████║╚██████╔╝██████╔╝███████╗██║ ╚████║╚██████╔╝██║ ╚═╝ ██║
        ╚══════╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝ ╚═════╝ ╚═╝     ╚═╝
`)
    fmt.Println(RESET)
    fmt.Println("                (c) Crafted with precision by =>> GAMER_0_0\n")
    fmt.Println("               Subenum v3 – Professional Recon Framework")
    fmt.Println("            For Advanced Security Researchers and Pentesters")
    fmt.Println(" *------------------------------------------------------------------------*")

    stop := make(chan bool)
    var wg sync.WaitGroup
    wg.Add(1)
    go spinner("Initializing Modules...", stop, &wg)
    time.Sleep(2 * time.Second)
    stop <- true
    wg.Wait()

    fmt.Println(GREEN + "   [+] Modules initialized successfully!\n" + RESET)
}

/* ================= HELP ================= */

func usage() {
    banner()
    fmt.Println(`Usage:
  subenum [options]

Target Options:
  -d string    Target domain (example.com | *.example.com | *.example.*)
  -l file      File with list of domains

Output Options:
  -o dir       Output directory (default: subdomain_enu)

General:
  -h           Show this help message

Examples:
  subenum -d example.com
  subenum -d "*.example.*"
  subenum -l scope.txt -o /path/recon
`)
    os.Exit(0)
}

/* ================= UTILS ================= */

// runTool يقوم بتشغيل أمر خارجي مع التحقق من وجوده
func runTool(name string, args ...string) []string {
    // التحقق من وجود الأداة في المسار (PATH)
    if _, err := exec.LookPath(name); err != nil {
        return []string{}
    }

    cmd := exec.Command(name, args...)
    out, err := cmd.Output()
    if err != nil {
        return []string{}
    }
    lines := strings.Split(string(out), "\n")
    var res []string
    for _, l := range lines {
        l = strings.TrimSpace(strings.ToLower(l))
        if l != "" {
            res = append(res, l)
        }
    }
    return res
}

func dedup(input []string) []string {
    seen := make(map[string]bool)
    var res []string
    for _, v := range input {
        if !seen[v] {
            seen[v] = true
            res = append(res, v)
        }
    }
    return res
}

func readLines(path string) []string {
    f, err := os.Open(path)
    if err != nil {
        return []string{}
    }
    defer f.Close()

    var lines []string
    sc := bufio.NewScanner(f)
    for sc.Scan() {
        t := strings.TrimSpace(sc.Text())
        if t != "" {
            lines = append(lines, strings.ToLower(t))
        }
    }
    return lines
}

func writeLines(path string, lines []string) {
    f, _ := os.Create(path)
    defer f.Close()
    for _, l := range lines {
        f.WriteString(l + "\n")
    }
}

/* ================= CRTSH (NATIVE GO IMPLEMENTATION) ================= */

type crtShEntry struct {
    NameValue string `json:"name_value"`
}

func getCrtSh(domain string) []string {
    apiURL := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", url.QueryEscape(domain))

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(apiURL)
    if err != nil {
        return []string{}
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return []string{}
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return []string{}
    }

    var data []crtShEntry
    if err := json.Unmarshal(body, &data); err != nil {
        return []string{}
    }

    seen := make(map[string]bool)
    var subs []string

    for _, entry := range data {
        lines := strings.Split(entry.NameValue, "\n")
        for _, line := range lines {
            sub := strings.TrimSpace(strings.ToLower(line))
            // إزالة wildcard إذا وجدت (*.domain.com)
            sub = strings.TrimPrefix(sub, "*.")

            if sub != "" && !seen[sub] && strings.Contains(sub, ".") { // تأكد بسيط من أنه نطاق
                seen[sub] = true
                subs = append(subs, sub)
            }
        }
    }
    return subs
}

/* ================= CHAOS ================= */

func chaosKey() string {
    cfg := filepath.Join(os.Getenv("HOME"), ".config/subenum")
    keyFile := filepath.Join(cfg, "chaos.key")

    // محاولة قراءة المفتاح الموجود
    if data, err := os.ReadFile(keyFile); err == nil {
        return strings.TrimSpace(string(data))
    }

    fmt.Print(YELLOW + "[?] Do you have Chaos API key? (y/n): " + RESET)
    in := bufio.NewReader(os.Stdin)
    ans, _ := in.ReadString('\n')
    ans = strings.TrimSpace(ans)

    if ans != "y" {
        return ""
    }

    fmt.Print(YELLOW + "[?] Enter Chaos API key: " + RESET)
    key, _ := in.ReadString('\n')
    key = strings.TrimSpace(key)

    os.MkdirAll(cfg, 0700)
    os.WriteFile(keyFile, []byte(key), 0600)
    return key
}

/* ================= MAIN ================= */

func main() {
    flag.StringVar(&domain, "d", "", "")
    flag.StringVar(&list, "l", "", "")
    flag.StringVar(&outDir, "o", "subdomain_enu", "")
    flag.Usage = usage
    flag.Parse()

    if domain == "" && list == "" {
        usage()
    }

    banner()

    var targets []string
    if domain != "" {
        targets = append(targets, domain)
    }
    if list != "" {
        targets = append(targets, readLines(list)...)
    }
    targets = dedup(targets)

    if len(targets) == 0 {
        fmt.Println(RED + "[-] No valid targets found." + RESET)
        os.Exit(1)
    }

    chaosAPI := chaosKey()
    in := bufio.NewReader(os.Stdin)

    for i, t := range targets {
        fmt.Printf(BLUE+"[+] Processing domain (%d/%d): %s\n"+RESET, i+1, len(targets), t)

        base := strings.TrimPrefix(t, "*.")
        outPath := filepath.Join(outDir, base)
        os.MkdirAll(outPath, 0755)

        var all []string
        var mu sync.Mutex
        var wg sync.WaitGroup

        // === 1. Subfinder (Goroutine) ===
        wg.Add(1)
        go func() {
            defer wg.Done()
            res := runTool("subfinder", "-all","-silent", "-d", base)
            mu.Lock()
            all = append(all, res...)
            mu.Unlock()
            fmt.Printf("    ├─ subfinder   : %d found\n", len(res))
        }()

        // === 2. Assetfinder (Goroutine) ===
        wg.Add(1)
        go func() {
            defer wg.Done()
            res := runTool("assetfinder", "--subs-only", base)
            mu.Lock()
            all = append(all, res...)
            mu.Unlock()
            fmt.Printf("    ├─ assetfinder : %d found\n", len(res))
        }()

        // === 3. Findomain (Goroutine) - NEW ADDITION ===
        wg.Add(1)
        go func() {
            defer wg.Done()
            // تشغيل findomain مع الوضع الصامت (-q)
            res := runTool("findomain", "-t", base, "-q")
            mu.Lock()
            all = append(all, res...)
            mu.Unlock()
            fmt.Printf("    ├─ findomain   : %d found\n", len(res))
        }()

        // === 4. Chaos (Goroutine) ===
        wg.Add(1)
        go func() {
            defer wg.Done()
            var res []string

            if chaosAPI != "" {
                res = runTool("chaos", "-d", base, "-key", chaosAPI)
                mu.Lock()
                all = append(all, res...)
                mu.Unlock()
                fmt.Printf("    ├─ chaos       : %d found\n", len(res))
            } else {
                // UX Improvement: Explicit message
                fmt.Printf("    ├─ chaos       : skipped (no API key)\n")
            }
        }()

        // === 5. Crt.sh (Goroutine) ===
        wg.Add(1)
        go func() {
            defer wg.Done()
            res := getCrtSh(base)
            mu.Lock()
            all = append(all, res...)
            mu.Unlock()
            fmt.Printf("    └─ crt.sh      : %d found\n", len(res))
        }()

        // انتظار انتهاء جميع الأدوات
        wg.Wait()

        all = dedup(all)
        fmt.Printf(GREEN+"    [+] Total unique subdomains: %d\n"+RESET, len(all))

        // كتابة الملف الأولي
        mainOutput := filepath.Join(outPath, "subdomains.txt")
        writeLines(mainOutput, all)

        fmt.Print(YELLOW + "[?] Do you want to add additional subdomains? (y/n): " + RESET)
        ans, _ := in.ReadString('\n')
        ans = strings.TrimSpace(ans)

        inputFile := mainOutput
        finalFileName := "subdomains.txt"

        if ans == "y" {
            fmt.Println(YELLOW + "[*] Enter subdomains (one per line). Empty line to finish:" + RESET)
            for {
                line, _ := in.ReadString('\n')
                line = strings.TrimSpace(strings.ToLower(line))
                if line == "" {
                    break
                }
                all = append(all, line)
            }
            all = dedup(all)
            finalFileName = "all_subdomains.txt"
            inputFile = filepath.Join(outPath, finalFileName)
            writeLines(inputFile, all)
        }

        // === HTTPX (Alive Check) مع Spinner ===
        httpxOut := filepath.Join(outPath, "httpx.txt")

        // إعداد وتشغيل الـ Spinner
        stopSpinner := make(chan bool)
        var wgSpinner sync.WaitGroup
        wgSpinner.Add(1)
        go spinner("Running httpx (alive check)...", stopSpinner, &wgSpinner)

        cmd := exec.Command("httpx", "-l", inputFile, "-silent", "-o", httpxOut)
        err := cmd.Run()

        // إيقاف الـ Spinner
        stopSpinner <- true
        wgSpinner.Wait()

        if err != nil {
            fmt.Println(RED + "    [-] Failed to run httpx. Is it installed?" + RESET)
        } else {
            fmt.Println(GREEN + "    [+] httpx completed successfully." + RESET)
        }

        // === Filtering Status Code 200 ===
        fmt.Println(BLUE + "[*] Filtering status code 200..." + RESET)
        httpx200Out := filepath.Join(outPath, "httpx_200.txt")
        cmd2 := exec.Command("httpx", "-l", httpxOut, "-mc", "200", "-silent", "-o", httpx200Out)
        if err := cmd2.Run(); err != nil {
            fmt.Println(RED + "    [-] Failed to filter httpx." + RESET)
        }

        fmt.Printf(GREEN+"[✔] Completed: %s (%d unique subdomains)\n\n"+RESET, t, len(all))
    }

    fmt.Println(GREEN + "[✔] Recon completed successfully." + RESET)
}
