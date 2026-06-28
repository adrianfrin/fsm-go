package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/flandersrin/workflow-go/workflow"
	"github.com/flandersrin/workflow-go/workflowtest"
)

const demoInstanceID = "order-demo-1"

var homePage = template.Must(template.New("home").Parse(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>订单流程 Demo</title>
  <style>
    body { margin: 0; background: #f7f4ee; color: #20242a; font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    main { width: min(760px, calc(100% - 32px)); margin: 32px auto; }
    h1 { margin: 0 0 8px; font-size: 30px; line-height: 1.15; letter-spacing: 0; }
    p { margin: 0; color: #667085; }
    form { display: flex; gap: 10px; margin: 24px 0; }
    input, button { border: 1px solid #d7d0c4; border-radius: 8px; font: inherit; }
    input { flex: 1; min-width: 0; padding: 11px 12px; background: #fffdf8; color: #20242a; }
    button { padding: 11px 16px; background: #1f766f; color: white; font-weight: 700; cursor: pointer; }
    ul { list-style: none; margin: 0; padding: 0; display: grid; gap: 10px; }
    li { display: flex; justify-content: space-between; gap: 12px; align-items: center; background: #fffdf8; border: 1px solid #d7d0c4; border-radius: 8px; padding: 12px 14px; }
    a { color: #1f766f; font-weight: 700; text-decoration: none; overflow-wrap: anywhere; }
    .hint { margin: 18px 0 10px; font-weight: 700; color: #20242a; }
    @media (max-width: 560px) {
      main { width: min(100% - 24px, 760px); margin: 20px auto; }
      form, li { align-items: stretch; flex-direction: column; }
    }
  </style>
</head>
<body>
  <main>
    <h1>订单流程 Demo</h1>
    <p>输入一个实例 ID，就能打开对应的订单流程时间线。</p>
    <form method="get" action="/">
      <input name="instance" placeholder="例如 order-demo-2" value="{{ .DefaultID }}" autocomplete="off">
      <button type="submit">打开时间线</button>
    </form>
    <p class="hint">已有实例</p>
    <ul>
      {{ range .Instances }}
      <li><a href="/instances/{{ .PathID }}/timeline">{{ .ID }}</a><a href="/api/instances/{{ .PathID }}/timeline">JSON</a></li>
      {{ end }}
    </ul>
  </main>
</body>
</html>`))

var timelinePage = template.Must(template.New("timeline").Funcs(template.FuncMap{
	"formatTime": func(t time.Time) string {
		if t.IsZero() {
			return "-"
		}
		return t.Local().Format("2006-01-02 15:04:05")
	},
	"formatPayload": func(payload map[string]any) string {
		if len(payload) == 0 {
			return "-"
		}
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return fmt.Sprint(payload)
		}
		return string(data)
	},
}).Parse(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>订单流程时间线</title>
  <style>
    body { margin: 0; background: #f7f4ee; color: #20242a; font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    main { width: min(960px, calc(100% - 32px)); margin: 32px auto; }
    header { display: flex; justify-content: space-between; gap: 20px; align-items: flex-end; border-bottom: 1px solid #d7d0c4; padding-bottom: 18px; }
    h1 { margin: 0 0 8px; font-size: 30px; line-height: 1.15; letter-spacing: 0; }
    p, time { margin: 0; color: #667085; }
    a { color: #1f766f; font-weight: 700; text-decoration: none; }
    .summary { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 12px; margin: 22px 0; }
    .box, .item { background: #fffdf8; border: 1px solid #d7d0c4; border-radius: 8px; }
    .box { padding: 14px; }
    .label { display: block; color: #667085; font-size: 13px; margin-bottom: 6px; }
    .value { display: block; font-size: 18px; font-weight: 700; overflow-wrap: anywhere; }
    ol { list-style: none; margin: 0; padding: 0; border-left: 2px solid #d7d0c4; }
    li { position: relative; margin-left: 18px; padding: 0 0 18px 22px; }
    li::before { content: ""; position: absolute; left: -27px; top: 6px; width: 12px; height: 12px; border-radius: 50%; background: #1f766f; box-shadow: 0 0 0 4px #f7f4ee; }
    .item { padding: 14px 16px; }
    .topline { display: flex; align-items: baseline; justify-content: space-between; gap: 16px; margin-bottom: 8px; }
    .message { font-weight: 700; overflow-wrap: anywhere; }
    .meta { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 10px; }
    .pill { border: 1px solid #d7d0c4; border-radius: 999px; color: #667085; font-size: 12px; padding: 3px 8px; background: #fbfaf6; }
    pre { margin: 0; white-space: pre-wrap; overflow-wrap: anywhere; color: #334155; font-size: 13px; line-height: 1.45; }
    @media (max-width: 720px) {
      main { width: min(100% - 24px, 960px); margin: 20px auto; }
      header, .topline { align-items: flex-start; flex-direction: column; }
      .summary { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <main>
    <header>
      <div>
        <h1>订单流程时间线</h1>
        <p>实例 {{ .Instance.ID }} 已跑完一次完整流程。</p>
      </div>
      <a href="{{ index .Links "json" }}">查看 JSON</a>
    </header>
    <section class="summary" aria-label="流程结果">
      <div class="box"><span class="label">当前状态</span><span class="value">{{ .Instance.State }}</span></div>
      <div class="box"><span class="label">运行结果</span><span class="value">{{ .Instance.Status }}</span></div>
      <div class="box"><span class="label">时间线条数</span><span class="value">{{ len .History }}</span></div>
    </section>
    <ol>
      {{ range .History }}
      <li>
        <div class="item">
          <div class="topline"><span class="message">{{ .Message }}</span><time>{{ formatTime .CreatedAt }}</time></div>
          <div class="meta">
            <span class="pill">{{ .Type }}</span>
            {{ if .State }}<span class="pill">状态 {{ .State }}</span>{{ end }}
            {{ if .Event }}<span class="pill">事件 {{ .Event }}</span>{{ end }}
            {{ if .Task }}<span class="pill">任务 {{ .Task }}</span>{{ end }}
          </div>
          <pre>{{ formatPayload .Payload }}</pre>
        </div>
      </li>
      {{ end }}
    </ol>
  </main>
</body>
</html>`))

type demoApp struct {
	runtime     *workflow.Runtime
	mu          sync.Mutex
	report      workflow.RunReport
	instanceIDs []string
	seen        map[string]bool
}

type homeView struct {
	DefaultID string
	Instances []instanceLink
}

type instanceLink struct {
	ID     string
	PathID string
}

type timelineView struct {
	Instance *workflow.WorkflowInstance  `json:"instance"`
	History  []workflow.ExecutionHistory `json:"history"`
	Report   workflow.RunReport          `json:"report"`
	Links    map[string]string           `json:"links"`
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	app, err := newDemoApp(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handleHome)
	mux.HandleFunc("/timeline", app.handleTimelineByQuery)
	mux.HandleFunc("/api/timeline", app.handleTimelineJSONByQuery)
	mux.HandleFunc("/instances/", app.handleTimelinePage)
	mux.HandleFunc("/api/instances/", app.handleTimelineJSON)

	url := displayURL(*addr, "/")
	fmt.Printf("demo started: %s\n", url)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func newDemoApp(ctx context.Context) (*demoApp, error) {
	store := workflowtest.NewMemoryStore()
	runtime := workflow.NewRuntime(store)
	Register(runtime)
	RegisterTask(runtime, HandlerPaymentCharge, workflow.TaskHandlerFunc(chargePayment))

	app := &demoApp{runtime: runtime, seen: map[string]bool{}}
	if err := app.ensureInstance(ctx, demoInstanceID); err != nil {
		return nil, err
	}
	return app, nil
}

func chargePayment(context.Context, workflow.TaskContext) (workflow.TaskResult, error) {
	return workflow.TaskResult{Output: map[string]any{"charged": true}}, nil
}

func (a *demoApp) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	instanceID := cleanInstanceID(r.URL.Query().Get("instance"))
	if instanceID != "" {
		http.Redirect(w, r, timelinePagePath(instanceID), http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := homePage.Execute(w, homeView{DefaultID: "order-demo-2", Instances: a.instanceLinks()}); err != nil {
		log.Printf("render home: %v", err)
	}
}

func (a *demoApp) handleTimelineByQuery(w http.ResponseWriter, r *http.Request) {
	instanceID := cleanInstanceID(r.URL.Query().Get("instance"))
	if instanceID == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, timelinePagePath(instanceID), http.StatusFound)
}

func (a *demoApp) handleTimelineJSONByQuery(w http.ResponseWriter, r *http.Request) {
	instanceID := cleanInstanceID(r.URL.Query().Get("instance"))
	if instanceID == "" {
		http.Error(w, "instance is required", http.StatusBadRequest)
		return
	}
	a.writeTimelineJSON(w, r, instanceID)
}

func (a *demoApp) handleTimelinePage(w http.ResponseWriter, r *http.Request) {
	instanceID, ok := timelineInstanceID(r.URL.Path, "/instances/")
	if !ok {
		http.NotFound(w, r)
		return
	}
	view, err := a.timelineView(r.Context(), instanceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := timelinePage.Execute(w, view); err != nil {
		log.Printf("render timeline: %v", err)
	}
}

func (a *demoApp) handleTimelineJSON(w http.ResponseWriter, r *http.Request) {
	instanceID, ok := timelineInstanceID(r.URL.Path, "/api/instances/")
	if !ok {
		http.NotFound(w, r)
		return
	}
	a.writeTimelineJSON(w, r, instanceID)
}

func (a *demoApp) writeTimelineJSON(w http.ResponseWriter, r *http.Request, instanceID string) {
	view, err := a.timelineView(r.Context(), instanceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(view); err != nil {
		log.Printf("write timeline json: %v", err)
	}
}

func (a *demoApp) ensureInstance(ctx context.Context, instanceID string) error {
	instanceID = cleanInstanceID(instanceID)
	if instanceID == "" {
		return fmt.Errorf("instance id is required")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.seen[instanceID] {
		return nil
	}

	if _, err := Start(ctx, a.runtime, instanceID, map[string]any{"amount": 100, "currency": "CNY"}); err != nil {
		return err
	}
	report, err := a.runtime.RunDueTasks(ctx, workflow.RunOptions{Limit: 10})
	if err != nil {
		return err
	}

	a.instanceIDs = append(a.instanceIDs, instanceID)
	a.seen[instanceID] = true
	a.report = report
	return nil
}

func (a *demoApp) instanceLinks() []instanceLink {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]instanceLink, 0, len(a.instanceIDs))
	for _, id := range a.instanceIDs {
		out = append(out, instanceLink{ID: id, PathID: url.PathEscape(id)})
	}
	return out
}

func (a *demoApp) currentReport() workflow.RunReport {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.report
}

func (a *demoApp) timelineView(ctx context.Context, instanceID string) (timelineView, error) {
	if err := a.ensureInstance(ctx, instanceID); err != nil {
		return timelineView{}, err
	}
	current, err := a.runtime.GetWorkflow(ctx, instanceID)
	if err != nil {
		return timelineView{}, err
	}
	history, err := a.runtime.ListHistory(ctx, instanceID)
	if err != nil {
		return timelineView{}, err
	}
	return timelineView{
		Instance: current,
		History:  history,
		Report:   a.currentReport(),
		Links: map[string]string{
			"page": timelinePagePath(instanceID),
			"json": timelineJSONPath(instanceID),
		},
	}, nil
}

func timelineInstanceID(path string, prefix string) (string, bool) {
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, "/timeline") {
		return "", false
	}
	instanceID := strings.TrimSuffix(strings.TrimPrefix(path, prefix), "/timeline")
	instanceID, err := url.PathUnescape(instanceID)
	if err != nil {
		return "", false
	}
	instanceID = cleanInstanceID(instanceID)
	if instanceID == "" {
		return "", false
	}
	return instanceID, true
}

func cleanInstanceID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || strings.Contains(value, "/") {
		return ""
	}
	return value
}

func timelinePagePath(instanceID string) string {
	return "/instances/" + url.PathEscape(instanceID) + "/timeline"
}

func timelineJSONPath(instanceID string) string {
	return "/api/instances/" + url.PathEscape(instanceID) + "/timeline"
}

func displayURL(addr string, path string) string {
	if strings.HasPrefix(addr, ":") {
		return "http://localhost" + addr + path
	}
	return "http://" + addr + path
}
