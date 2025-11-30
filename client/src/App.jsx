import './App.css'
import { BrowserRouter as Router, Routes, Route, Navigate, BrowserRouter } from 'react-router-dom'
import Home from './pages/home/Home'
import Details from './pages/details/Details'
import Navbar from './pages/navbar/Navbar'
import { SocketProvider } from './contexts/SocketContext'

export default function App() {
  
  return (
    <SocketProvider>
      <BrowserRouter>
        <Navbar />
        <Routes>
          <Route path='/' element={<Navigate to='/home' replace />} />
          <Route path='/home' element={<Home />} />
          <Route path='/details/:id' element={<Details />} />
        </Routes>
      </BrowserRouter>
    </SocketProvider>
  )
}
