import React, { Component } from 'react';
import { Layout, Menu, Icon } from 'antd';
import './style.scss';

const SubMenu = Menu.SubMenu;
const Header = Layout.Header;

export default class Navigation extends Component {

  render() {

    return (
      <div>
        <Header style={{ padding: '0px', height: '49px' }}>
          <div className="logo">
	    sern
	  </div>
          <Menu mode="horizontal" defaultSelectedKeys={['1']} >
            <Menu.Item key="1">Dashboard</Menu.Item>
            <Menu.Item key="2">Documentation</Menu.Item>
            <SubMenu
              style={{ float: 'right' }}
              title={<span><span>Walid </span><Icon type="user" /></span>}
            >
              <Menu.Item>
	        <a href="/logout">Logout</a>
	      </Menu.Item>
            </SubMenu>
          </Menu>
        </Header >
      </div >
    );
  }
}
