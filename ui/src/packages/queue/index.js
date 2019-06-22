import React, { Component } from 'react';
import { Table, Tag } from 'antd';


function getStatus(status) {
  let color = 'green';
  if (status === 'paused') {
    color = 'volcano';
  }
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

const data = [
  {
    key: '1',
    repo: 'fakerr/experiment2',
    active: '#4674',
    queue: '#4674 #46546 #7897 #1325',
    status: 'wait', // either running, paused or wait(pending)
  }
];

export default class PRTable extends Component {

  render() {

    return (
	<Table columns={columns} dataSource={data} pagination={false} style={{ backgroundColor: '#FFFFFF' }} />
    );
  }
}
