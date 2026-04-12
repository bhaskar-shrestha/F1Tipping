import React, { useState, useEffect } from 'react';
import API from '../api';

function UserPage() {
  const [selectedDrivers, setSelectedDrivers] = useState([]);
  const [selectedTeams, setSelectedTeams] = useState([]);
  const [driverList, setDriverList] = useState([]);
  const [teamList, setTeamList] = useState([]);
  const [predictions, setPredictions] = useState([]);
  const [myPredictions, setMyPredictions] = useState([]);

  useEffect(() => {
    loadDriverList();
    loadTeamList();
    loadMyPredictions();
  }, []);

  const loadDriverList = () => {
    API.get('/api/admin/drivers').then(setDriverList);
  };

  const loadTeamList = () => {
    API.get('/api/admin/teams').then(setTeamList);
  };

  const loadMyPredictions = () => {
    API.get('predictions/user/my-user').then(setMyPredictions);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    API.post('/api/predictions', {
      user_id: 'my-user',
      driver_ids: selectedDrivers.map(d => d.id),
      team_ids: selectedTeams.map(t => t.id),
    }).then(() => {
      alert('Prediction submitted!');
      loadMyPredictions();
    });
  };

  const getResultColor = (points) => {
    if (points >= 40) return { color: 'green', fontWeight: 'bold' };
    if (points >= 20) return { color: 'orange' };
    return { color: 'gray' };
  };

  return (
    <div>
      <h2>User View - Select Your Predictions</h2>

      {/* Driver Selection */}
      <div style={{ marginTop: '20px' }}>
        <h3>Select 5 Drivers</h3>
        <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '20px' }}>
          <thead>
            <tr style={{ borderBottom: '2px solid #ccc' }}>
              <th style={{ padding: '10px' }}>Select</th>
              <th style={{ padding: '10px' }}>Name</th>
              <th style={{ padding: '10px' }}>Constructor</th>
            </tr>
          </thead>
          <tbody>
            {driverList.map((driver) => (
              <tr key={driver.id} style={{ borderBottom: '1px solid #eee' }}>
                <td style={{ padding: '10px' }}>
                  <input
                    type="checkbox"
                    checked={selectedDrivers.some(d => d.id === driver.id)}
                    onChange={() => {
                      if (selectedDrivers.length < 5) {
                        setSelectedDrivers([...selectedDrivers, driver]);
                      }
                    }}
                  />
                </td>
                <td style={{ padding: '10px' }}>{driver.name}</td>
                <td style={{ padding: '10px' }}>{driver.constructor_name}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Team Selection */}
      <div style={{ marginTop: '20px' }}>
        <h3>Select 2 Teams</h3>
        <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '20px' }}>
          <thead>
            <tr style={{ borderBottom: '2px solid #ccc' }}>
              <th style={{ padding: '10px' }}>Select</th>
              <th style={{ padding: '10px' }}>Constructor Name</th>
            </tr>
          </thead>
          <tbody>
            {teamList.map((team) => (
              <tr key={team.id} style={{ borderBottom: '1px solid #eee' }}>
                <td style={{ padding: '10px' }}>
                  <input
                    type="checkbox"
                    checked={selectedTeams.some(t => t.id === team.id)}
                    onChange={() => {
                      if (selectedTeams.length < 2) {
                        setSelectedTeams([...selectedTeams, team]);
                      }
                    }}
                  />
                </td>
                <td style={{ padding: '10px' }}>{team.constructorName}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Submit Button */}
      <button
        onClick={handleSubmit}
        disabled={selectedDrivers.length !== 5 || selectedTeams.length !== 2}
        style={{
          padding: '10px 20px',
          fontSize: '16px',
          backgroundColor: selectedDrivers.length === 5 && selectedTeams.length === 2 ? '#1976d2' : '#ccc',
          color: 'white',
          border: 'none',
          borderRadius: '4px',
          cursor: selectedDrivers.length === 5 && selectedTeams.length === 2 ? 'pointer' : 'not-allowed',
          marginTop: '20px',
        }}
      >
        Submit Prediction
      </button>

      {/* My Predictions Results */}
      <div style={{ marginTop: '30px' }}>
        <h3>My Predictions</h3>
        {myPredictions.length === 0 ? (
          <p style={{ color: 'gray' }}>No predictions yet</p>
        ) : (
          myPredictions.map((pred) => {
            const resultColor = getResultColor(pred.totalPoints);
            return (
              <div key={pred.id} style={{
                padding: '15px',
                border: '1px solid #ddd',
                borderRadius: '8px',
                marginBottom: '10px',
                backgroundColor: '#f9f9f9',
                ...resultColor,
              }}>
                <div><strong>Drivers:</strong> {pred.driver_ids?.join(', ') || 'N/A'}</div>
                <div><strong>Teams:</strong> {pred.team_ids?.join(', ') || 'N/A'}</div>
                <div><strong>Sprint Points:</strong> {pred.sprint_points || 0}</div>
                <div><strong>Race Points:</strong> {pred.race_points || 0}</div>
                <div><strong>Total Points:</strong> {pred.total_points || 0}</div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}

export default UserPage;
