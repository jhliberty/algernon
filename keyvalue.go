package main

import (
	"net/http"

	"github.com/xyproto/permissions2"
	"github.com/xyproto/simpleredis"
	"github.com/yuin/gopher-lua"
)

// Identifier for the Set class in Lua
const lKeyValueClass = "KEYVALUE"

// Get the first argument, "self", and cast it from userdata to a set.
func checkKeyValue(L *lua.LState) *simpleredis.KeyValue {
	ud := L.CheckUserData(1)
	if kv, ok := ud.Value.(*simpleredis.KeyValue); ok {
		return kv
	}
	L.ArgError(1, "keyvalue expected")
	return nil
}

// Create a new KeyValue collection.
// id is the name of the KeyValue colleciton.
// dbindex is the Redis database index (typically 0).
func newKeyValue(L *lua.LState, pool *simpleredis.ConnectionPool, id string, dbindex int) (*lua.LUserData, error) {
	// Create a new simpleredis set
	kv := simpleredis.NewKeyValue(pool, id)
	kv.SelectDatabase(dbindex)
	// Create a new userdata struct
	ud := L.NewUserData()
	ud.Value = kv
	L.SetMetatable(ud, L.GetTypeMetatable(lKeyValueClass))
	return ud, nil
}

// String representation
// Returns the name of the KeyValue collection
// tostring(kv) -> string
func kvToString(L *lua.LState) int {
	L.Push(lua.LString("keyvalue"))
	return 1 // Number of returned values
}

// Set a key and value. Returns true if it worked out.
// kv:set(string, string) -> bool
func kvSet(L *lua.LState) int {
	kv := checkKeyValue(L) // arg 1
	key := L.ToString(2)
	value := L.ToString(3)
	L.Push(lua.LBool(nil == kv.Set(key, value)))
	return 1 // Number of returned values
}

// Takes a key, returns a value. May return an empty string.
// kv:get(string) -> string
func kvGet(L *lua.LState) int {
	kv := checkKeyValue(L) // arg 1
	key := L.ToString(2)
	retval, err := kv.Get(key)
	if err != nil {
		retval = ""
	}
	L.Push(lua.LString(retval))
	return 1 // Number of returned values
}

// Remove a key. Returns true if it worked out.
// kv:del(string) -> bool
func kvDel(L *lua.LState) int {
	kv := checkKeyValue(L) // arg 1
	value := L.ToString(2)
	L.Push(lua.LBool(nil == kv.Del(value)))
	return 1 // Number of returned values
}

// Remove the keyvalue itself. Returns true if it worked out.
// kv:remove() -> bool
func kvRemove(L *lua.LState) int {
	kv := checkKeyValue(L) // arg 1
	L.Push(lua.LBool(nil == kv.Remove()))
	return 1 // Number of returned values
}

// The keyvalue methods that are to be registered
var kvMethods = map[string]lua.LGFunction{
	"__tostring": kvToString,
	"set":        kvSet,
	"get":        kvGet,
	"del":        kvDel,
	"remove":     kvRemove,
}

// Make functions related to HTTP requests and responses available to Lua scripts
func exportKeyValue(w http.ResponseWriter, req *http.Request, L *lua.LState, userstate *permissions.UserState) {
	pool := userstate.Pool()
	dbindex := userstate.DatabaseIndex()

	// Register the KeyValue class and the methods that belongs with it.
	mt := L.NewTypeMetatable(lKeyValueClass)
	mt.RawSetH(lua.LString("__index"), mt)
	L.SetFuncs(mt, kvMethods)

	// The constructor for new KeyValues takes a name and an optional redis db index
	L.SetGlobal("KeyValue", L.NewFunction(func(L *lua.LState) int {
		name := L.ToString(1)

		// Check if the optional argument is given
		localDBIndex := dbindex
		if L.GetTop() == 2 {
			localDBIndex = L.ToInt(2)
		}

		// Create a new keyvalue in Lua
		userdata, err := newKeyValue(L, pool, name, localDBIndex)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			L.Push(lua.LNumber(1))
			return 3 // Number of returned values
		}

		// Return the keyvalue object
		L.Push(userdata)
		return 1 // Number of returned values
	}))
}
