
import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import Dashboard from './pages/Dashboard'
import NotFound from './pages/NotFound'
import Auth from './pages/Auth'

const App: React.FunctionComponent = (): JSX.Element => {
  return <BrowserRouter>
    <Switch>
      <Route exact={true} path="/" component={Dashboard} />
      <Route exact={true} path="/auth" component={Auth} />
      <Route path='*' exact={true} component={NotFound} />
    </Switch>
  </BrowserRouter>
}

export default App;