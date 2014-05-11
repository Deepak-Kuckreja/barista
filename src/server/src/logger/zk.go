package main

import "fmt"
import "launchpad.net/gozk/zookeeper"

type ZK struct {
  Conn *zookeeper.Conn
}

func Make() (*ZK, error) {
  zk := &ZK{}
  conn, session, err := zookeeper.Dial("localhost:2181", 5e9)
  if err != nil {
    fmt.Printf("Can't connect to zookeeper: %v\n", err)
    return nil, err
  }
  defer conn.Close()

  // Wait for connection.
  event := <-session
  if event.State != zookeeper.STATE_CONNECTED {
    fmt.Printf("Can't connect to zookeeper: %v\n", event)
    return nil, err
  }
  zk.Conn = conn
  return zk, nil
}

func (zk *ZK) Write(key string, value string) error {
  stats, err := zk.Conn.Exists(key)
  if err != nil {
      fmt.Println("Error creating or writing to file")
      return err
  }
  if stats != nil {
    _, err = zk.Conn.Create(key, value, 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
    if err != nil {
      fmt.Printf("Can't create key: %v\n", err)
      return err
    }
  }
  stats, err = zk.Conn.Set(key, value, -1)
  return err
}

func (zk *ZK) Read(key string) (string, error) {
  data, stats, err := zk.Conn.Get(key)
  return data, err
}


func main() {
  zk, err := Make()
  zk.Write("K", "V")
  data, err := zk.Read("K")
  fmt.Println(data)
}