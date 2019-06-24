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
      repo: 'Fakerr/experiment2',
      loading: true,
      data: []
    };
  }

  refreshQueue() {
    this.getData();
  }

  componentDidMount() {
    this.getData();
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
      });
  };

  render() {

    return (
	<Table loading={this.state.loading} columns={columns} dataSource={this.state.data} pagination={false} style={{ backgroundColor: '#FFFFFF' }} />
    );
  }
}
