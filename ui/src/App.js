import React, { Component } from 'react';
import { Collapse, Layout } from 'antd';
import Navigation from './packages/navigation';
import Queue from './packages/queue';
import 'antd/dist/antd.css';
import './App.scss';

const { Panel } = Collapse;
const { Content, Sider } = Layout;

function callback(key) {
  console.log(key);
}

const text = `
  A dog is a type of domesticated animal.
  Known for its loyalty and faithfulness,
  it can be found as a welcome guest in many households across the world.
`;

const customPanelStyle = {
  background: '#f7f7f7',
  color: 'green'
};

const node = (
    <p style={{ color: '#1890ff', textDecoration: 'underline', marginBottom: '0' }}>My Repositories</p>
);


export default class App extends Component {
  render() {
    return (
      <div className="wrapper">
	<div className="box header">
          <Navigation></Navigation>
        </div> 
        <div className="box content">
          <Layout>
            <Sider width={230} style={{ background: '#f0f2f5' }}>
              <Collapse defaultActiveKey={['1']} onChange={callback} style={{ borderRadius: '0' }}>
                <Panel header={node} key="1" style={customPanelStyle}>
                  <p>{text}</p>
                </Panel>
              </Collapse>
            </Sider>
            <Layout style={{ padding: '0 5px 10px' }}>
              <Content>
                <Queue></Queue>
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
