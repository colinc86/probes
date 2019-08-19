# probes
The probes package helps examine numerical values that change over time by providing methods to collect values safely across goroutines through the `Probe` structure.

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

p.Flush()
p.Deactivate()
```

### Examining the Input Signal
Examine the probe's input signal at any time by using the `Signal` and `RecentValue` methods.

```go
// "9"
fmt.Println(p.RecentValue())

// "[0 1 2 3 4 5 6 7 8 9]"
fmt.Println(p.Signal())
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

### Blocking and Non-Blocking
`Probe` provides two ways to pass input to the probe:
- The input channel, `C`. (non-blocking)
- The `Push` method. (blocking)

Use whichever is more appropriate, but be aware that calling `Deactivate` on a probe does not automatically flush the probes input channel buffer. You must call `Flush` before calling `Deactivate` to ensure that all values will be represented in the probe's signal.

Mixing input methods by using both the channel and push method is supported by the push method's `flush` parameter. Set this to true to preserve input order in the probe's signal.

```go
p := NewProbe()
p.Activate()

// Some work...
for var i := 0; i < 10; i++ {
  p.C <- float64(i)
}

p.Push(10, true)
p.Deactivate()
```
