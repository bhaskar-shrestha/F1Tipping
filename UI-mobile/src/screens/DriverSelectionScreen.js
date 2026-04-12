import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, Alert, Dimensions } from 'react-native';
import { Link } from 'react-navigation';

const { width } = Dimensions.get('window');

export default function DriverSelectionScreen({ navigation }) {
  const [drivers, setDrivers] = useState([]);
  const [selectedDrivers, setSelectedDrivers] = useState([]);

  useEffect(() => {
    // Load drivers from API
    fetchDrivers();
  }, []);

  const fetchDrivers = async () => {
    try {
      const API_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(`${API_URL}/api/admin/drivers`);
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      const data = await response.json();
      setDrivers(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Error loading drivers:', error);
      // Fallback mock data
      setDrivers([
        { id: 'd1', name: 'Max Verstappen', constructor_name: 'Red Bull' },
        { id: 'd2', name: 'Sergio Perez', constructor_name: 'Red Bull' },
        { id: 'd3', name: 'Charles Leclerc', constructor_name: 'Ferrari' },
        { id: 'd4', name: 'Carlos Sainz', constructor_name: 'Ferrari' },
        { id: 'd5', name: 'Lewis Hamilton', constructor_name: 'Mercedes' },
        { id: 'd6', name: 'George Russell', constructor_name: 'Mercedes' },
        { id: 'd7', name: 'Lando Norris', constructor_name: 'McLaren' },
        { id: 'd8', name: 'Oscar Piastri', constructor_name: 'McLaren' },
      ]);
    }
  };

  const toggleDriver = (driverId) => {
    if (selectedDrivers.length < 5) {
      if (selectedDrivers.includes(driverId)) {
        setSelectedDrivers(selectedDrivers.filter(id => id !== driverId));
      } else {
        setSelectedDrivers([...selectedDrivers, driverId]);
      }
    } else {
      Alert.alert('Limit Reached', 'You can only select 5 drivers');
    }
  };

  const submitSelection = () => {
    if (selectedDrivers.length === 5) {
      navigation.navigate('TeamSelection');
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.header}>Select 5 Drivers</Text>

      <ScrollView style={styles.driverList}>
        {drivers.map(driver => (
          <TouchableOpacity
            key={driver.id}
            style={[styles.driverItem, selectedDrivers.includes(driver.id) && styles.selectedItem]}
            onPress={() => toggleDriver(driver.id)}
          >
            <Text style={styles.driverName}>{driver.name}</Text>
            <Text style={styles.constructorName}>{driver.constructor_name}</Text>
            {selectedDrivers.includes(driver.id) && (
              <Text style={styles.checkmark}>✓</Text>
            )}
          </TouchableOpacity>
        ))}
      </ScrollView>

      <TouchableOpacity style={styles.submitButton} onPress={submitSelection}>
        <Text style={styles.submitButtonText}>Next</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
    backgroundColor: '#fff',
  },
  header: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  driverList: {
    flex: 1,
    gap: 8,
  },
  driverItem: {
    padding: 12,
    borderRadius: 8,
    backgroundColor: '#f5f5f5',
    flexDirection: 'row',
    alignItems: 'center',
  },
  selectedItem: {
    backgroundColor: '#1976d2',
  },
  driverName: {
    fontSize: 16,
    fontWeight: '600',
  },
  constructorName: {
    fontSize: 14,
    color: '#666',
  },
  checkmark: {
    color: '#fff',
    fontSize: 24,
    marginLeft: 'auto',
  },
  submitButton: {
    padding: 16,
    backgroundColor: '#1976d2',
    borderRadius: 8,
    alignItems: 'center',
  },
  submitButtonText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: '600',
  },
});
