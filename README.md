#  check_status

**check_status** is a Poller using Rest API. You can configure the environment as well.

## Quick Start / Installation

To get started with **check_status**, you'll need to do the following:

- Install the latest version of Go.
- Install the **check_status** package using go get.
- Use the **check_status** package in your Go code.
- Here's an example of how to install the package:

```
go get github.com/lanyere/check_status
```

## Examples

```Go
    p := NewPoller(&Config{
        Providers: []string{"specify all your providers"},
        Interval:  2, // set the interval or use default
        Log:       yourLogger,
        Database:  yourDB,
    })
    
    // start the polling
    p.StartPolling()
    
    // get the status of your desired transaction
    status, err = p.GetTransactionStatus("ID of your transaction")
    // handle errors and status
    
    // don't forget to stop the polling when you're done with polling
    p.StopPolling
```

## Interfaces

**check_status** includes an interfaces, which allows you to use your **own tools**.
This means you can create your own logger / database and use them with **check_status** package

Here's the link to [Interfaces](https://github.com/lanyere/check_status/blob/main/interfaces.go).

```
https://github.com/lanyere/check_status/interfaces.go
```

## Ecosystem / Documentations

- [net/http for REST API](https://pkg.go.dev/net/http)
- [time.Ticker for intervals](https://pkg.go.dev/time)
- [prometheus for metrics](https://github.com/prometheus/client_golang)


## Contributors

> **check_status** is maintained by [lanyere](https://github.com/lanyere).
> If you'd like to contribute to the project, please reach out to us via [Issues](https://github.com/lanyere/check_status/issues).