import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ProtectedRoute } from './components/ProtectedRoute';
import Layout from './components/Layout/Layout';
import Landing from './pages/Landing/Landing';
import Login from './pages/Login/Login';
import WorkerDashboard from './pages/Worker/Dashboard';
import PlannedOrders from './pages/Worker/PlannedOrders';
import OrderHistory from './pages/Worker/OrderHistory';
import OrderDetail from './pages/Worker/OrderDetail';
import OwnerDashboard from './pages/Owner/Dashboard';
import OwnerRequests from './pages/Owner/Requests';
import OwnerPlannedOrders from './pages/Owner/PlannedOrders';
import OwnerOrderHistory from './pages/Owner/OrderHistory';
import OwnerOrderDetail from './pages/Owner/OrderDetail';
import OwnerCreateOrder from './pages/Owner/CreateOrder';
import OwnerStatistics from './pages/Owner/Statistics';
import OwnerWorkers from './pages/Owner/Workers';

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/login" element={<Login />} />

          {/* Worker */}
          <Route path="/dashboard" element={
            <ProtectedRoute role="worker"><Layout><WorkerDashboard /></Layout></ProtectedRoute>
          } />
          <Route path="/orders/planned" element={
            <ProtectedRoute role="worker"><Layout><PlannedOrders /></Layout></ProtectedRoute>
          } />
          <Route path="/orders/history" element={
            <ProtectedRoute role="worker"><Layout><OrderHistory /></Layout></ProtectedRoute>
          } />
          <Route path="/orders/:id" element={
            <ProtectedRoute role="worker"><Layout><OrderDetail /></Layout></ProtectedRoute>
          } />

          {/* Owner */}
          <Route path="/owner/dashboard" element={
            <ProtectedRoute role="owner"><Layout><OwnerDashboard /></Layout></ProtectedRoute>
          } />
          <Route path="/owner/requests" element={
            <ProtectedRoute role="owner"><Layout><OwnerRequests /></Layout></ProtectedRoute>
          } />
          <Route path="/owner/orders/planned" element={
            <ProtectedRoute role="owner"><Layout><OwnerPlannedOrders /></Layout></ProtectedRoute>
          } />
          <Route path="/owner/orders/history" element={
            <ProtectedRoute role="owner"><Layout><OwnerOrderHistory /></Layout></ProtectedRoute>
          } />
          <Route path="/owner/orders/:id" element={
            <ProtectedRoute role="owner"><Layout><OwnerOrderDetail /></Layout></ProtectedRoute>
          } />
          <Route path="/owner/orders/create" element={
            <ProtectedRoute role="owner"><Layout><OwnerCreateOrder /></Layout></ProtectedRoute>
          } />
          <Route path="/owner/statistics" element={
            <ProtectedRoute role="owner"><Layout><OwnerStatistics /></Layout></ProtectedRoute>
          } />
          <Route path="/owner/workers" element={
            <ProtectedRoute role="owner"><Layout><OwnerWorkers /></Layout></ProtectedRoute>
          } />

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}