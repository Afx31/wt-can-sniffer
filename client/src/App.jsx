import './App.css'
import { BrowserRouter as Router, Routes, Route, Navigate, useParams } from 'react-router-dom'
import Home from './pages/home/Home'
import Details from './pages/details/Details'
import { SocketProvider } from './contexts/SocketContext'

export default function App() {
  const { id } = useParams()
  
  return (
    <SocketProvider>
      <Router>
        <Routes>
          <Route path='/' element={<Navigate to='/home' replace />} />
          <Route path='/home' element={<Home />} />
          <Route path='/details/:id' element={<Details />} />
        </Routes>
      </Router>
    </SocketProvider>
  )
}
