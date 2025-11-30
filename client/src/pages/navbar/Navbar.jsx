import './Navbar.css'
import { NavLink } from 'react-router-dom'
import { useSocket } from '../../contexts/SocketContext'

export default function Navbar() {
	const { isConnected, handleConnection, handleDisconnection, sendReqToServer } = useSocket()

	function handleSaveToFile() {
		sendReqToServer(2, 'test').then((res) => {
		})
	}

	return (
		<nav className='navbar-container'>
			<div className='navbar-left'>
				<NavLink to='/' className='navbar-item'>Home</NavLink>
				<button
					className='btn-save-to-file'
					onClick={handleSaveToFile}
				>
					Save to file
				</button>
					
			</div>
			<div className='navbar-right'>
				<button
					className='btn-server-connect'
					onClick={isConnected ? handleDisconnection : handleConnection}
				>
					{isConnected ? 'Disconnect from server' : 'Connect to server'}
				</button>
				<p
					className={`p-server-connect ${isConnected ? 'p-server-connected' : ''}`}
				>
					{isConnected ? 'Connected' : 'Disconnected'}
				</p>
			</div>
		</nav>
	)
}