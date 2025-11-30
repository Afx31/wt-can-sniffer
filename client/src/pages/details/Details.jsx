import './Details.css'
import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { useSocket } from '../../contexts/SocketContext'

export default function Details() {
	const { id } = useParams()
	const { socketRef, isConnected, fetchDataFromServer } = useSocket()
	const [canData, setCanData] = useState({})

	useEffect(() => {
		if (!isConnected || !socketRef.current) return

		fetchDataFromServer('GET-MSG-DATA', id).then((res) => {
			setCanData(res)
		})		
	}, [])

	return (
		<div className='d-container'>
			<h2>CAN ID: { id }</h2>
			<table>
				<thead>
				<tr>
					<th>Index</th>
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
					{canData?.data?.reverse().map((item, index) => (
						<tr key={canData.data.length - 1 - index}>
							<td>{canData.data.length - 1 - index}</td>
							{Object.values(item).map((v, i) => (
								<td key={i}>{v}</td>
							))}
						</tr>
					))}
				</tbody>
			</table>
		</div>
	)
}