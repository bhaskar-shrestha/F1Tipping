import React from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';

export default function ResultsScreen() {
  const [results, setResults] = React.useState([]);

  React.useEffect(() => {
    fetchResults();
  }, []);

  const fetchResults = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/predictions/user/my-user');
      const data = await response.json();
      setResults(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Error loading results:', error);
      // Fallback mock data
      setResults([
        {
          id: 'r1',
          driver_ids: ['d1', 'd3', 'd5', 'd7', 'd9'],
          team_ids: ['t1', 't2'],
          sprint_points: 25,
          race_points: 60,
          total_points: 85,
        },
      ]);
    }
  };

  const getPointsColor = (points) => {
    if (points >= 50) return '#4caf50';
    if (points >= 30) return '#ff9800';
    return '#f5f5f5';
  };

  return (
    <View style={styles.container}>
      <Text style={styles.header}>Your Results</Text>

      {results.length === 0 ? (
        <Text style={styles.noResults}>No results yet</Text>
      ) : (
        results.map(result => (
          <View key={result.id} style={styles.card}>
            <Text style={styles.sectionTitle}>Prediction</Text>
            <Text style={styles.driversText}>
              Drivers: {result.driver_ids?.join(', ') || 'N/A'}
            </Text>
            <Text style={styles.teamsText}>
              Teams: {result.team_ids?.join(', ') || 'N/A'}
            </Text>

            <Text style={styles.sectionTitle}>Points</Text>
            <View style={styles.pointsRow}>
              <View style={styles.pointsBox}>
                <Text style={styles.pointsLabel}>Sprint</Text>
                <Text style={[styles.pointsValue, { color: getPointsColor(result.sprint_points) }]}>
                  {result.sprint_points || 0}
                </Text>
              </View>
              <View style={styles.pointsBox}>
                <Text style={styles.pointsLabel}>Race</Text>
                <Text style={[styles.pointsValue, { color: getPointsColor(result.race_points) }]}>
                  {result.race_points || 0}
                </Text>
              </View>
            </View>

            <View style={styles.totalPointsBox}>
              <Text style={styles.totalPointsLabel}>Total Points</Text>
              <Text style={[styles.totalPointsValue, { color: getPointsColor(result.total_points) }]}>
                {result.total_points || 0}
              </Text>
            </View>
          </View>
        ))
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
    backgroundColor: '#f5f5f5',
  },
  header: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  noResults: {
    fontSize: 16,
    color: '#666',
    textAlign: 'center',
    marginTop: 40,
  },
  card: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 8,
  },
  driversText: {
    fontSize: 14,
    color: '#333',
    marginBottom: 4,
  },
  teamsText: {
    fontSize: 14,
    color: '#666',
    marginBottom: 12,
  },
  sectionTitlePoints: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 8,
  },
  pointsRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: 12,
  },
  pointsBox: {
    alignItems: 'center',
    padding: 8,
    borderRadius: 8,
    backgroundColor: '#f5f5f5',
  },
  pointsLabel: {
    fontSize: 12,
    color: '#666',
  },
  pointsValue: {
    fontSize: 24,
    fontWeight: 'bold',
  },
  totalPointsBox: {
    backgroundColor: '#1976d2',
    padding: 16,
    borderRadius: 8,
    alignItems: 'center',
  },
  totalPointsLabel: {
    color: '#fff',
    fontSize: 14,
  },
  totalPointsValue: {
    color: '#fff',
    fontSize: 32,
    fontWeight: 'bold',
    marginTop: 4,
  },
});
