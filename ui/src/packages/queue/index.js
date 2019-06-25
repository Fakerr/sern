import React, { Component } from 'react';
import { Table, Tag } from 'antd';


function getStatus(status) {
  let color = 'yellow';
  //if (status === 'paused') {
  //  color = 'volcano';
  //}
  if (status === 'running') {
    color = 'green';
  }
  if (status === 'wait') {
    color = 'yellow';
  }
  return (
    <Tag color={color} key={status}>
      {status.toUpperCase()}
    </Tag>
  )
}

const columns = [
  {
    title: 'Repository',
    dataIndex: 'repo',
    key: 'repo',
  },
  {
    title: 'Status',
    key: 'status',
    dataIndex: 'status',
    render: status => (
      <span>
        {getStatus(status)}
      </span>
    ),
  },
  {
    title: 'Active',
    dataIndex: 'active',
    key: 'active',
  },
  {
    title: 'Queued PRs',
    dataIndex: 'queue',
    key: 'queue',
  }
];

export default class PRTable extends Component {

  constructor(props) {
    super(props);

    this.state = {
      repo: this.props.repo,
      loading: false,
      data: []
    };
  }

  refreshQueue() {
    this.getData();
  }

  componentDidMount() {
    if (this.state.repo) {
      this.getData();
    }
  }

  componentWillReceiveProps({repo}) {
    this.setState({...this.state,repo}, () => {
      this.getData(); //why not use componentDidUpdate ?
    });
  }

  getData = () => {
    this.setState({ loading: true })
    fetch('/api/' + this.state.repo + '/queue')
      .then(response => response.json())
      .then(data => {
	const res = {
	  key: '1',
	  repo: this.state.repo,
	  status: 'wait'
	}
	if (data) {
	  res.active = '#' + data.Active;
	  res.status = data.Status;
	  res.queue = data.Queue.map(el => {
	    return '#' + el;
	  }).join(' ');
	  
	}
	this.setState({ data: [res], loading: false })
      }).catch(err =>{
	this.setState({ loading: false });
	alert("Internal error");
      });
  };

  render() {

    return (
	<Table loading={this.state.loading} columns={columns} dataSource={this.state.data} pagination={false} style={{ backgroundColor: '#FFFFFF' }} />
    );
  }
}
