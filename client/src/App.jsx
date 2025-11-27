import React from 'react';
import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const ws = new WebSocket("ws://localhost:8080/ws")
  const [canData, setCanData] = useState({})

  useEffect(() => {
    ws.onopen = () => {
      console.log('Connected')
      ws.send('Hello from client')
    }

    // ws.onmessage = (e) => {
      // console.log('Received from server: ', e.data)
    // }

    ws.onerror = (err) => {
      console.log('ERROR: ', err)
    }

    ws.addEventListener('message', handleWebSocketData)

    return () => ws.close()
  }, [])

  function handleWebSocketData(e) {
    const json = JSON.parse(e.data)
    // console.log(json)

    // Solution 1
    // setCanData(prev => ({
    //   ...prev,
    //   [json.id]: {
    //     id: json.id,
    //     length: json.length,
    //     data: [...json.data]  // new array
    //   }
    // }))

    // Solution 2
    setCanData(prev => ({
      ...prev,
      [json.id]: { ...json, data: [...json.data] }
    }))
  }

  return (
    <div className='container'>
      <button
        onClick={() => handleWebSocketConnection()}
      >
        Connection
      </button>
      <p>CAN Bus data</p>

      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Length</th>
            <th>Byte 0</th>
            <th>Byte 1</th>
            <th>Byte 2</th>
            <th>Byte 3</th>
            <th>Byte 4</th>
            <th>Byte 5</th>
            <th>Byte 6</th>
            <th>Byte 7</th>
          </tr>
        </thead>
        <tbody>
          {Object.values(canData).map(frame => (
            <tr key={frame.id}>
              <td>{frame.id}</td>
              <td>{frame.length}</td>
              {frame.data.map((d, i) => (
                <td key={i}>{d}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default App
