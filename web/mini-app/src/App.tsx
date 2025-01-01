import { useState } from 'react'

import './App.css'

import '@telegram-apps/telegram-ui/dist/styles.css'
import { AppRoot, Cell, List, Section, Tabbar } from '@telegram-apps/telegram-ui'

import WebApp from '@twa-dev/sdk'

const cellsTexts = ['Chat Settings', 'Data and Storage', 'Devices'];

function App() {
  const [count] = useState(0)
  const tabs = [{ id: 1, text: "Карта" }, { id: 1, text: "Events" }, { id: 1, text: "Stars" }]
  const [currentTab, setCurrentTab] = useState(tabs[0].id);

  return (
    <AppRoot>
      {/* List component to display a collection of items */}
      <List>
        {/* Section component to group items within the list */}
        <Section header="Header for the section" footer="Footer for the section">
          {/* Mapping through the cells data to render Cell components */}
          {cellsTexts.map((cellText, index) => (
            <Cell key={index}>
              {cellText}
            </Cell>
          ))}
        </Section>
        <div className="card">
          <button onClick={() => WebApp.showAlert(`Hello World! Current count is ${count}`)}>
            Show Alert
          </button>
        </div>
      </List>
      <a href='/not_found'>go to iditusi website</a>
      <Tabbar >
        {tabs.map(({
          id,
          text,
        }) => <Tabbar.Item key={id} text={text} selected={id === currentTab} onClick={() => setCurrentTab(id)}>
            <div className='' style={{ height: 40 }}></div>
          </Tabbar.Item>)}
      </Tabbar>
    </AppRoot >
  )
}

export default App