package bluetoothemulation

// Code generated by cdproto-gen. DO NOT EDIT.

// EventGattOperationReceived event for when a GATT operation of |type| to
// the peripheral with |address| happened.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/BluetoothEmulation#event-gattOperationReceived
type EventGattOperationReceived struct {
	Address string            `json:"address"`
	Type    GATTOperationType `json:"type"`
}
