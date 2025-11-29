import './Home.css'
import { useSocket } from '../../contexts/SocketContext';
import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'

// import React from 'react';

export default function Home() {
  const { socketRef, isConnected, handleConnection, handleDisconnection } = useSocket()
  const [canData, setCanData] = useState({})

  useEffect(() => {
    if (!isConnected || !socketRef.current) return

    const socket = socketRef.current
    socket.addEventListener('message', handleWebSocketData)

    //return () => ws.close()
    return () => socket.removeEventListener('message', handleWebSocketData)
  }, [isConnected])

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
    <div className='h-container'>
      <button
        onClick={isConnected ? handleDisconnection : handleConnection}
      >
        {isConnected ? 'Disconnect' : 'Connect'}
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
              <td>
                <Link to={`/details/${frame.id}`}>{frame.id}</Link>
              </td>
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
