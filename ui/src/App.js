import React, { Component } from 'react';
import { Tabs, Layout, Button } from 'antd';
import Navigation from './packages/navigation';
import RepoList from './packages/repositories';
import PRTable from './packages/queue';
import 'antd/dist/antd.css';
import './App.scss';

const { Content, Sider } = Layout;
const { TabPane } = Tabs;

function callback(key) {
  console.log(key);
}

export default class App extends Component {

  constructor(props) {
    super(props);
    this.child = React.createRef();

    this.state = {
      title: '',
      selectedItem: null
    };
  }

  refresh = () => {
    this.child.current.refreshQueue();
  }

  selectItem = (ev) => {
    this.setState({ selectedItem: ev.currentTarget.value, title: ev.currentTarget.value });
    console.log('Selected repository: ' + ev.currentTarget.value);
  }

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
              <RepoList selectItem={this.selectItem}></RepoList>
            </Sider>
            <Layout style={{ padding: '0 20px 10px' }}>
              <Content>
                <h1 className="repo-name">
	          {this.state.title}
                </h1>
	        <Tabs defaultActiveKey="1"
		      onChange={callback}
		      tabBarExtraContent={<Button disabled={this.state.selectedItem == null}
						  onClick={this.refresh}>Refresh Queue</Button>}>
	          <TabPane tab="Merge Queue" key="1">
                    <PRTable repo={this.state.selectedItem} ref={this.child}></PRTable>
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
