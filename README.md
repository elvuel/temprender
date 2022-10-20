# Temprender

Temprender package is for `GO` template render and rendering transfer, part of my `NEP`(Nesting Ecological Program) infrastructure.

## Intro

Frequent or high similarity(repeated) `light complexity` scaffolding steps are required in NEP.

Up an unified base processing to reduces some variety and repeat, while simplify expand and management.

## Usage

```shell
go get github.com/elvuel/temprender
```

## Quickstart

* checkout examples/quick

```go
    // ...
    const tmpl = "quick.tmpl"
    t := task.NewTask() 
    t.PerformCtx, _ = context.NewDefaultContext()
    t.PerformCtx.S("user", "Quick")
    t.PerformCtx.S("greeter", "temprender") 

    t.LoadTempatesFromMap(map[string]string{
    tmpl: `Hi {{.G "user"}}, 
    Greeting from {{.G "greeter"}}!`,
    }) 
    transporter, _ := puppet.NewPuppeteer()
    t.RegisterTransporters("default", transporter) 
    t.RenderAll(stdctx.TODO()) 
    t.Transport(stdctx.TODO(), "default")
    // ...
```
