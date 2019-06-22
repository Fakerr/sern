import React, { Component } from 'react';
import { Collapse, List, Typography } from 'antd';

const { Panel } = Collapse;

const text = `
 No repository found
`;

const node = (
    <p style={{ color: '#1890ff', textDecoration: 'underline', marginBottom: '0' }}>My Repositories</p>
);

const data = [
  'fakerr/experiment2',
  'fakerr/sern'
];

const customPanelStyle = {
  background: '#f7f7f7',
  color: 'green'
};

export default class RepoList extends Component {

  render() {

    return (
      <div>
        <Collapse defaultActiveKey={['1']} style={{ borderRadius: '0' }}>
          <Panel header={node} key="1" style={customPanelStyle}>
	    <List
              bordered
              dataSource={data}
              renderItem={item => (
		<List.Item>
		  <Typography.Text strong>{item}</Typography.Text>
		</List.Item>
	      )}
	    />
          </Panel>
        </Collapse>
      </div >
    );
  }
}
