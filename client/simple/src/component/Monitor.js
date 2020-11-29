import React, { Component } from 'react';

class Monitor extends Component {
    constructor(props) {
        super(props)
        this.state = {
            drone : null,
            server : null,
        }
    }

    componentDidMount() {
        this.time = setInterval(
          () => this.tick(),
          1000
        );
    }

    componentWillUnmount() {
        clearInterval(this.timer);
    }

    tick() {
        fetch("http://3.91.1.79:10000/state")
        .then( res => res.json())
        .then( 
            (result) => {
                console.log(result)
		result.drone.dji.altitude = result.drone.dji.gps[2]
                this.setState({
                    drone : {
                        dji : result.drone.dji,
                        saturatr : result.drone.saturatr,
			signal: result.drone.signal, 
			upload : result.drone.upload,
			download : result.drone.download,
			ping: result.drone.ping,
			iperf: result.drone.iperf,

                    },
                    server : {
                        saturatr : result.server.saturatr,
			iperf: result.server.iperf,
                    }
                });
            },
            (error) => {
                console.error("Problem loading state")
            }
        )
    }

    render() {
        return (
            <div>
                <div>
                    <h2>Drone</h2>
                    { this.state.drone && 
                        <div>
                            <Dji status={this.state.drone.dji}></Dji>
			    <Signal status={this.state.drone.signal}></Signal>
			    <Iperf status={this.state.drone.iperf}></Iperf>
			    <Ping status={this.state.drone.ping}></Ping>
			    <p>upload: {this.state.drone.upload} Mbps</p>
			    <p>download: {this.state.drone.download} Mbps</p>
                        </div>
                    }
                </div>
                <div>
                    <h2>Server</h2>
                    { this.state.server && 
			<div>
			    <Iperf status={this.state.server.iperf}></Iperf>
			</div>
                    }
                </div>
            </div>
        )
    }
}

export default Monitor

class Iperf extends Component {
	render() {
	    return ( 
	        <div>
		    <h3>Iperf</h3>
		    <p>bandwidth: {this.props.status.bandwidth}{this.props.status.unit}</p>
  	        </div>
	    )
	}

}

class Ping extends Component {
	render() {
	    return (
	        <div>
		    <h3>Ping</h3>
		    <p>latency: {this.props.status.latency}ms </p>
  	        </div>
	    )
	}
}

class Signal extends Component {
    render() {
	return (
	    <div>
		<h3>Hilink</h3>
		<p>rsrp: {this.props.status.rsrp} </p>
		<p>rsrq: {this.props.status.rsrq} </p>
		<p>rssi: {this.props.status.rssi} </p>
		<p>sinr: {this.props.status.sinr} </p>
		<p>cell_id: {this.props.status.cell_id} </p>
	    </div>

	)
    }

}

class Dji extends Component {
    render() {
        return (
            <div>
                <h3>Dji</h3>
                <p>battery: {this.props.status.battery} </p>
                <p>altitude: {this.props.status.altitude} </p>
            </div>
        )
    }
}

class Saturatr extends Component {
    render() {
        return (
            <div>
                <h3>Saturatr</h3>
                <div>
                    <h4>Acker</h4>
                    <p>sent: {this.props.status.acker.sent}</p>
                    <p>received: {this.props.status.acker.received}</p>
                </div>
                <div>
                    <h4>Saturatr</h4>
                    <p>sent: {this.props.status.saturatr.sent}</p>
                    <p>received: {this.props.status.saturatr.received}</p>
                </div>
            </div>
        )
    }
}
