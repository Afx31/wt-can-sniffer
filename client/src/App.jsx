import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const ws = new WebSocket("ws://localhost:8080/ws")
  const [canData, setCanData] = useState([])

  useEffect(() => {
    ws.onopen = () => {
      console.log('Connected')
      ws.send('Hello from client')
    }

    ws.onmessage = (e) => {
      console.log('Received from server: ', e.data)
    }

    ws.onerror = (err) => {
      console.log('ERROR: ', err)
    }

    ws.addEventListener('message', handleWebSocketData)

    return () => ws.close()
  }, [])

  function handleWebSocketData(e) {
    const data = JSON.parse(e.data)
    // console.log(data)
    setCanData(prev => [...prev, data])
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
            <th>Data</th>
          </tr>
        </thead>
        <tbody>
          {canData.map((data, i) => (
            <tr
              key={data.id}
            >
              <td>{data.id}</td>
              <td>{data.length}</td>
              <td>{data.data}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default App
