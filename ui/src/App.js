import React, { Component } from 'react';
import { Tabs, Collapse, Layout, List, Button, Typography } from 'antd';
import Navigation from './packages/navigation';
import Queue from './packages/queue';
import 'antd/dist/antd.css';
import './App.scss';

const { Panel } = Collapse;
const { Content, Sider } = Layout;
const { TabPane } = Tabs;

function callback(key) {
  console.log(key);
}

const text = `
 No repository found
`;

const customPanelStyle = {
  background: '#f7f7f7',
  color: 'green'
};

const node = (
    <p style={{ color: '#1890ff', textDecoration: 'underline', marginBottom: '0' }}>My Repositories</p>
);

const data = [
  'fakerr/experiment2',
  'fakerr/sern'
];


export default class App extends Component {
  render() {
    return (
      <div className="wrapper">
	<div className="box header">
          <Navigation></Navigation>
        </div> 
        <div className="box content">
          <Button style={{ margin: '7px'}} type="dashed">Add / Remove repositories</Button>
          <Layout>
            <Sider width={260} style={{ background: '#f0f2f5' }}>
              <Collapse defaultActiveKey={['1']} onChange={callback} style={{ borderRadius: '0' }}>
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
            </Sider>
            <Layout style={{ padding: '0 20px 10px' }}>
              <Content>
	        <Tabs defaultActiveKey="1" onChange={callback}>
	          <TabPane tab="Merge Queue" key="1">
                    <Queue></Queue>
                  </TabPane>
                  <TabPane tab="Settings" key="2">
	            Content of settings
                  </TabPane>
	        </Tabs>
              </Content>
            </Layout>
          </Layout>
        </div>
        <div className="box footer">
          <p><a href="http://predictix.com">&copy; Infor 2018</a></p>
        </div>
      </div>
    );
  }
}
