# Go-CLI

This is is a minimal CLI library for Go. It is designed to get out of your way and let you focus on your application rather than trying to figure out how to do X with the CLI library.

There are no external dependencies and basically no frills and that is the main feature of this library.
It does not try to do argument parsing, instead you can use one of the many flag libraries that already exist.

## Usage


Basic example:

```go
cmd := NewCmd("ping", func(ctx context.Context) error {
    fmt.Println("pong")
    return nil
})

_ = cmd.Run(context.Background(), os.Args[1:])
```

You can add flags to your command:

```go
var extra string
cmd := NewCmd("ping", func(ctx context.Context) error {
    fmt.Println("pong", extra)
    return nil
})

cmd.Flags().StringVar(&extra, "extra", "", "extra string")
_ = cmd.Run(context.Background(), os.Args[1:])
```

You can add subcommands to your command:

```go
cmd := NewCmd("top-level", nil)

cmd.NewCmd("sub", func(ctx context.Context) error {
    fmt.Println("sub")
    return nil
})

_ = cmd.Run(context.Background(), os.Args[1:])
```

You can add subcommands to subcommands, add flags to subcommands, etc.

By default this will use the built-in go flag library, but you can use any flag library you want by implementing the `FlagSet` interface.

```go
cmd := NewCmdWithFlagSet("ping", func(ctx context.Context) error {
    fmt.Println("pong")
    return nil
    // function will create a new flag set using the pflag library from github.com/spf13/pflag
}, func(name string) *pflag.FlagSet { return pflag.NewFlagSet(name, pflag.ExitOnError) })
```
