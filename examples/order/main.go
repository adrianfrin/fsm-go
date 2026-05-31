package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/flandersrin/workflow-go/workflow"
	"github.com/flandersrin/workflow-go/workflowtest"
)

const demoInstanceID = "order-demo-1"

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
      <a href="/api/instances/{{ .Instance.ID }}/timeline">查看 JSON</a>
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
	runtime *workflow.Runtime
	report  workflow.RunReport
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
	mux.HandleFunc("/", app.redirectToTimeline)
	mux.HandleFunc("/instances/", app.handleTimelinePage)
	mux.HandleFunc("/api/instances/", app.handleTimelineJSON)

	url := displayURL(*addr, "/instances/"+demoInstanceID+"/timeline")
	fmt.Printf("demo started: %s\n", url)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func newDemoApp(ctx context.Context) (*demoApp, error) {
	store := workflowtest.NewMemoryStore()
	runtime := workflow.NewRuntime(store)
	Register(runtime)
	RegisterTask(runtime, HandlerPaymentCharge, workflow.TaskHandlerFunc(chargePayment))

	_, err := Start(ctx, runtime, demoInstanceID, map[string]any{"amount": 100, "currency": "CNY"})
	if err != nil {
		return nil, err
	}
	report, err := runtime.RunDueTasks(ctx, workflow.RunOptions{Limit: 10})
	if err != nil {
		return nil, err
	}
	return &demoApp{runtime: runtime, report: report}, nil
}

func chargePayment(context.Context, workflow.TaskContext) (workflow.TaskResult, error) {
	return workflow.TaskResult{Output: map[string]any{"charged": true}}, nil
}

func (a *demoApp) redirectToTimeline(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/instances/"+demoInstanceID+"/timeline", http.StatusFound)
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

func (a *demoApp) timelineView(ctx context.Context, instanceID string) (timelineView, error) {
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
		Report:   a.report,
		Links: map[string]string{
			"page": "/instances/" + instanceID + "/timeline",
			"json": "/api/instances/" + instanceID + "/timeline",
		},
	}, nil
}

func timelineInstanceID(path string, prefix string) (string, bool) {
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, "/timeline") {
		return "", false
	}
	instanceID := strings.TrimSuffix(strings.TrimPrefix(path, prefix), "/timeline")
	if instanceID == "" || strings.Contains(instanceID, "/") {
		return "", false
	}
	return instanceID, true
}

func displayURL(addr string, path string) string {
	if strings.HasPrefix(addr, ":") {
		return "http://localhost" + addr + path
	}
	return "http://" + addr + path
}
