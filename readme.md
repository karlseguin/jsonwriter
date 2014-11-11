# JsonWriter

Write JSON to an io.Writer. Useful for simple cases where you want to avoid encoding/json or require greater control.

## Usage

JsonWriter can serialize strings, ints (unsigned/signed, 8/16/32/64), floats, bools and nulls values.


```go
buffer := new(bytes.Buffer)
writer := jsonWriter.New(buffer)

writer.RootObject(func(){
  writer.KeyValue("id", "leto")
  writer.KeyValue("age", 3000)
  writer.KeyValue("god", true)
  writer.Object("sister", func(){
    writer.KeyValue("id", "ghanima")
    writer.KeyValue("alive", false)
  })
  writer.Array("friends", func(){
    writer.Value("duncan", "moneo")
  })
})
```

* `RootObject(nest func())` - Generate a document as an object
* `RootArray(nest func())` - Generate a document as an array
* `KeyValue(key string, value interface)` - Write the key value pair
* `Value(value interface)` - Write the value (only useful within an array)
* `Object(key string, nest func())` - Start a nested object
* `Array(key string, nest func())` - Start an array
* `ArrayObject(nest func())` - Used to place an object within an array

# Misc

The [Typed](https://github.com/karlseguin/typed) library provides the opposite functionality: a lightweight approaching to reading JSON data.
