import React, { useState, useEffect } from 'react';
import API from '../../api';

function AdminPage() {
  const [drivers, setDrivers] = useState([]);
  const [teams, setTeams] = useState([]);
  const [driverForm, setDriverForm] = useState({ name: '', constructor_id: '', constructor_name: '' });
  const [teamForm, setTeamForm] = useState({ constructorName: '' });

  useEffect(() => {
    loadDrivers();
    loadTeams();
  }, []);

  const loadDrivers = () => {
    API.get('/api/admin/drivers').then(setDrivers);
  };

  const loadTeams = () => {
    API.get('/api/admin/teams').then(setTeams);
  };

  const handleSubmitDriver = (e) => {
    e.preventDefault();
    API.post('/api/admin/drivers', driverForm).then(() => {
      setDriverForm({ name: '', constructor_id: '', constructor_name: '' });
      loadDrivers();
      loadTeams();
    });
  };

  const handleSubmitTeam = (e) => {
    e.preventDefault();
    API.post('/api/admin/teams', teamForm).then(() => {
      setTeamForm({ constructorName: '' });
      loadTeams();
    });
  };

  return (
    <div>
      <h2>Admin Page</h2>

      <div style={{ display: 'flex', gap: '40px', marginTop: '20px' }}>
        {/* Add Driver Form */}
        <div style={{ flex: 1 }}>
          <h3>Add Driver</h3>
          <form onSubmit={handleSubmitDriver} style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <input
              type="text"
              placeholder="Driver Name"
              value={driverForm.name}
              onChange={(e) => setDriverForm({ ...driverForm, name: e.target.value })}
              required
            />
            <input
              type="text"
              placeholder="Constructor ID"
              value={driverForm.constructor_id}
              onChange={(e) => setDriverForm({ ...driverForm, constructor_id: e.target.value })}
              required
            />
            <input
              type="text"
              placeholder="Constructor Name"
              value={driverForm.constructor_name}
              onChange={(e) => setDriverForm({ ...driverForm, constructor_name: e.target.value })}
              required
            />
            <button type="submit">Add Driver</button>
          </form>
        </div>

        {/* Add Team Form */}
        <div style={{ flex: 1 }}>
          <h3>Add Team</h3>
          <form onSubmit={handleSubmitTeam} style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <input
              type="text"
              placeholder="Constructor Name"
              value={teamForm.constructorName}
              onChange={(e) => setTeamForm({ constructorName: e.target.value })}
              required
            />
            <button type="submit">Add Team</button>
          </form>
        </div>
      </div>

      {/* Driver List */}
      <div style={{ marginTop: '30px' }}>
        <h3>All Drivers</h3>
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ borderBottom: '2px solid #ccc' }}>
              <th style={{ padding: '10px' }}>Name</th>
              <th style={{ padding: '10px' }}>Constructor ID</th>
            </tr>
          </thead>
          <tbody>
            {drivers.map((driver) => (
              <tr key={driver.id} style={{ borderBottom: '1px solid #eee' }}>
                <td style={{ padding: '10px' }}>{driver.name}</td>
                <td style={{ padding: '10px' }}>{driver.constructor_id}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Team List */}
      <div style={{ marginTop: '30px' }}>
        <h3>All Teams</h3>
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ borderBottom: '2px solid #ccc' }}>
              <th style={{ padding: '10px' }}>Constructor Name</th>
              <th style={{ padding: '10px' }}>Constructor ID</th>
            </tr>
          </thead>
          <tbody>
            {teams.map((team) => (
              <tr key={team.id} style={{ borderBottom: '1px solid #eee' }}>
                <td style={{ padding: '10px' }}>{team.constructor_name}</td>
                <td style={{ padding: '10px' }}>{team.id}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default AdminPage;
