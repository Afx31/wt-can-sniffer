import { createContext, useContext, useRef, useState } from 'react';

const SocketContext = createContext(null)

export function SocketProvider({ children }) {
    const socketRef = useRef(null)
    const [isConnected, setIsConnected] = useState(false)

    function handleConnection() {
        if (socketRef.current) {
            console.log('[DEBUG] Already connected')
        }

        console.log('[DEBUG] Connecting to Server via WebSocket...')
        const ws = new WebSocket('ws://localhost:8080/ws')
        
        ws.onopen = () => {
            console.log('[DEBUG] Connected to Server')
            setIsConnected(true)
        }

        ws.onclose = () => {
            console.log('[DEBUG] WebSocket closed from Server')
            setIsConnected(false)
        }

        ws.onerror = (err) => console.log('[DEBUG] [ERROR] WebSocket error: ', err)

        socketRef.current = ws
    }

    function handleDisconnection() {
        socketRef.current?.close(1000, 'Client request disconnection')
        console.log('[DEBUG] Disconnecting from Server')
    }

    function fetchDataFromServer(type, param) {
        return new Promise((resolve) => {
            if (type === 'GET-MSG-DATA') {
                const reqObj = {
                    id: 1,
                    data: param
                }
                socketRef.current.send(JSON.stringify(reqObj))
                socketRef.current.onmessage = (e) => {
                    resolve(JSON.parse(e.data))
                }
            }
        })
    }

    return (
        <SocketContext.Provider value ={{ socketRef, isConnected, handleConnection, handleDisconnection, fetchDataFromServer}}>
            {children}
        </SocketContext.Provider>
    )
}

export function useSocket() {
    const context = useContext(SocketContext)
    if (!context) throw new Error('[ERROR] useSocket must be used within a SocketProvider')
    return context
}