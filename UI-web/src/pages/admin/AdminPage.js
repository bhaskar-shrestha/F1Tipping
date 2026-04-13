import React, { useState, useEffect } from 'react';
import API from '../../api';

function AdminPage() {
  const [drivers, setDrivers] = useState([]);
  const [teams, setTeams] = useState([]);
  const [driverForm, setDriverForm] = useState({ name: '', constructor_id: '', constructor_name: '' });
  const [teamForm, setTeamForm] = useState({ constructor_name: '' });
  const [driverError, setDriverError] = useState('');
  const [teamError, setTeamError] = useState('');
  const [driverSuccess, setDriverSuccess] = useState('');
  const [teamSuccess, setTeamSuccess] = useState('');

  useEffect(() => {
    loadDrivers();
    loadTeams();
  }, []);

  const loadDrivers = () => {
    API.get('/admin/drivers')
      .then(data => setDrivers(Array.isArray(data) ? data : []))
      .catch(err => {
        console.error('Failed to load drivers:', err);
        setDrivers([]);
      });
  };

  const loadTeams = () => {
    API.get('/admin/teams')
      .then(data => setTeams(Array.isArray(data) ? data : []))
      .catch(err => {
        console.error('Failed to load teams:', err);
        setTeams([]);
      });
  };

  const handleSubmitDriver = (e) => {
    e.preventDefault();
    setDriverError('');
    setDriverSuccess('');

    API.post('/admin/drivers', driverForm)
      .then(() => {
        setDriverForm({ name: '', constructor_id: '', constructor_name: '' });
        setDriverSuccess('Driver added successfully!');
        loadDrivers();
        loadTeams();
        // Clear success message after 3 seconds
        setTimeout(() => setDriverSuccess(''), 3000);
      })
      .catch(err => {
        console.error('Failed to add driver:', err);
        if (err.status === 400) {
          setDriverError(`Validation error: ${err.data || err.message}`);
        } else if (err.status === 404) {
          setDriverError('Endpoint not found.');
        } else {
          setDriverError(`Error adding driver: ${err.message}`);
        }
      });
  };

  const handleSubmitTeam = (e) => {
    e.preventDefault();
    setTeamError('');
    setTeamSuccess('');

    API.post('/admin/teams', teamForm)
      .then(() => {
        setTeamForm({ constructor_name: '' });
        setTeamSuccess('Team added successfully!');
        loadTeams();
        // Clear success message after 3 seconds
        setTimeout(() => setTeamSuccess(''), 3000);
      })
      .catch(err => {
        console.error('Failed to add team:', err);
        if (err.status === 400) {
          setTeamError(`Validation error: ${err.data || err.message}`);
        } else if (err.status === 404) {
          setTeamError('Endpoint not found.');
        } else {
          setTeamError(`Error adding team: ${err.message}`);
        }
      });
  };

  const messageStyle = (type) => ({
    padding: '12px',
    marginBottom: '12px',
    borderRadius: '4px',
    border: `1px solid ${type === 'error' ? '#ef5350' : '#81c784'}`,
    backgroundColor: type === 'error' ? '#ffebee' : '#e8f5e9',
    color: type === 'error' ? '#c62828' : '#2e7d32',
  });

  return (
    <div>
      <h2>Admin Page</h2>

      <div style={{ display: 'flex', gap: '40px', marginTop: '20px' }}>
        {/* Add Driver Form */}
        <div style={{ flex: 1 }}>
          <h3>Add Driver</h3>
          {driverError && <div style={messageStyle('error')}>{driverError}</div>}
          {driverSuccess && <div style={messageStyle('success')}>{driverSuccess}</div>}
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
          {teamError && <div style={messageStyle('error')}>{teamError}</div>}
          {teamSuccess && <div style={messageStyle('success')}>{teamSuccess}</div>}
          <form onSubmit={handleSubmitTeam} style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
            <input
              type="text"
              placeholder="Constructor Name"
              value={teamForm.constructor_name}
              onChange={(e) => setTeamForm({ constructor_name: e.target.value })}
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
