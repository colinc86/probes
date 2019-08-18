# probes
The probes package helps examine numerical values that change over time.

## Usage
### Activation/Deactivation
Activate a probe by calling its `Activate` method, and deactivate it by calling `Deactivate`. While doing work, send values to the probe's input channel.

```go
p := NewProbe()
p.Activate()

// Some work...
for var i := 0; i < 10; i++ {
  p.C <- float64(i)
}

p.Deactivate()
```

### Examining the Input Signal
Examine the probe's input signal at any time by using the `Signal` and `RecentValue` methods, or save a plot of the signal to file.

```go
// "9"
fmt.Println(p.RecentValue())

// "{0 1 2 3 4 5 6 7 8 9}"
fmt.Println(p.Signal())

// Plot the signal and save to testSignal.png
p.WriteSignalToPNG("testSignal")
```

### Reuse
Calling `Activate` _after_ a call to `Deactivate` will continue appending values to the probes previous input signal. To prevent this behavior, call the probe's `ClearSignal` method before a call to `Activate`.

```go
// ... some work
p.Deactivate()

// Log our input signal
fmt.Println(p.Signal())

// Clear the input and re-activate
p.ClearSignal()
p.Activate()
```

### Buffers and State
Set the probe's input channel buffer length by setting the `InputBufferLength` property. This value is set to 1 by default and will only take effect after the next call to `Activate`.

Set the probe's maximum signal length by setting the `MaximumSignalLength` property. This value is set to `math.MaxInt32` by default and should only be set while the probe is in an inactive state.

Check the state of the probe at any time with its `IsActive` method.

```go
p.InputBufferLength = 10

if !p.IsActive() {
  p.MaximumSignalLength = 100
}
```
