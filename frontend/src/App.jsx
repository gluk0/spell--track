import React, { useState, useEffect } from 'react';
import { 
  Container, 
  Paper, 
  Typography, 
  Box,
  CircularProgress,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Button,
  Alert
} from '@mui/material';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const API_BASE_URL = '/api';

function App() {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(false);
  const [caseId, setCaseId] = useState('');
  const [error, setError] = useState(null);
  const [availableCases, setAvailableCases] = useState([]);
  const [processedData, setProcessedData] = useState([]);

  const fetchEvents = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await fetch(`${API_BASE_URL}/events`);
      if (!response.ok) {
        throw new Error('Failed to fetch events');
      }
      const data = await response.json();
      console.log('Fetched events:', data);
      setEvents(data);
      
      // Extract unique case IDs
      const uniqueCases = [...new Set(data.map(event => event.case_id))];
      console.log('Available case IDs:', uniqueCases);
      setAvailableCases(uniqueCases.sort());
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const calculateEventDurations = () => {
    if (!caseId) return [];

    // Filter events for the selected case
    const caseEvents = events.filter(event => event.case_id === caseId);
    console.log('Events for case', caseId, ':', caseEvents);
    
    if (caseEvents.length === 0) {
      console.log('No events found for case:', caseId);
      return [];
    }
    
    // Group events by event_name
    const eventGroups = {};
    
    caseEvents.forEach(event => {
      if (!eventGroups[event.event_name]) {
        eventGroups[event.event_name] = {
          name: event.event_name,
          startEvents: [],
          endEvents: []
        };
      }

      try {
        // Parse the timestamp string to a Date object
        const timestamp = new Date(event.created_at);
        if (isNaN(timestamp.getTime())) {
          console.error(`Invalid timestamp for event:`, event);
          return;
        }
        
        if (event.event_type === 'start') {
          eventGroups[event.event_name].startEvents.push(timestamp.getTime());
        } else if (event.event_type === 'end') {
          eventGroups[event.event_name].endEvents.push(timestamp.getTime());
        }
      } catch (err) {
        console.error(`Error parsing timestamp for event:`, event, err);
      }
    });

    console.log('Event groups:', eventGroups);

    // Calculate durations for each event type
    const durations = Object.entries(eventGroups).map(([eventName, group]) => {
      // Sort timestamps to ensure proper matching
      group.startEvents.sort((a, b) => a - b);
      group.endEvents.sort((a, b) => a - b);

      let totalDuration = 0;
      let eventCount = 0;

      // Match start and end events
      for (let i = 0; i < Math.min(group.startEvents.length, group.endEvents.length); i++) {
        const startTime = group.startEvents[i];
        const endTime = group.endEvents[i];
        const duration = (endTime - startTime) / (1000 * 60); // Convert to minutes
        
        console.log(`Event ${eventName}:`, {
          startTime: new Date(startTime).toISOString(),
          endTime: new Date(endTime).toISOString(),
          duration: duration
        });

        if (duration > 0) {
          totalDuration += duration;
          eventCount++;
        }
      }

      const averageDuration = eventCount > 0 ? totalDuration / eventCount : 0;

      const result = {
        name: eventName,
        duration: Number(averageDuration.toFixed(2)),
        totalDuration: Number(totalDuration.toFixed(2)),
        count: eventCount
      };

      console.log(`Result for ${eventName}:`, result);
      return result;
    });

    console.log('Calculated durations:', durations);

    // Filter out events with no valid durations
    const filteredDurations = durations.filter(event => event.duration > 0);
    console.log('Filtered durations:', filteredDurations);
    return filteredDurations;
  };

  useEffect(() => {
    fetchEvents();
  }, []);

  useEffect(() => {
    if (caseId && events.length > 0) {
      console.log('Processing data for case:', caseId);
      const data = calculateEventDurations();
      console.log('Setting processed data:', data);
      setProcessedData(data);
    }
  }, [caseId, events]);

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Paper elevation={3} sx={{ p: 3 }}>
        <Typography variant="h4" gutterBottom>
          Case Event Durations
        </Typography>
        
        <Box sx={{ mb: 3, display: 'flex', gap: 2, alignItems: 'center' }}>
          <FormControl sx={{ minWidth: 200 }}>
            <InputLabel id="case-select-label">Case ID</InputLabel>
            <Select
              labelId="case-select-label"
              id="case-select"
              value={caseId}
              label="Case ID"
              onChange={(e) => setCaseId(e.target.value)}
              disabled={loading}
            >
              <MenuItem value="">
                <em>Select a case</em>
              </MenuItem>
              {availableCases.map((id) => (
                <MenuItem key={id} value={id}>
                  {id}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <Button 
            variant="contained" 
            onClick={() => fetchEvents()}
            disabled={loading}
          >
            {loading ? <CircularProgress size={24} /> : 'Refresh Data'}
          </Button>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {caseId && processedData.length === 0 && (
          <Alert severity="info" sx={{ mb: 2 }}>
            No events found for this case ID
          </Alert>
        )}

        {processedData.length > 0 && (
          <Box sx={{ height: 400 }}>
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={processedData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis 
                  dataKey="name" 
                  angle={-45}
                  textAnchor="end"
                  height={100}
                  interval={0}
                />
                <YAxis 
                  label={{ value: 'Duration (minutes)', angle: -90, position: 'insideLeft' }}
                  domain={[0, 'dataMax + 0.1']}
                />
                <Tooltip 
                  formatter={(value, name, props) => [
                    `${value} minutes (Total: ${props.payload.totalDuration} minutes, Count: ${props.payload.count})`,
                    'Average Duration'
                  ]}
                  labelFormatter={(label) => `Event Type: ${label}`}
                />
                <Legend />
                <Bar 
                  dataKey="duration" 
                  fill="#8884d8" 
                  name="Average Duration (minutes)"
                  maxBarSize={50}
                />
              </BarChart>
            </ResponsiveContainer>
          </Box>
        )}
      </Paper>
    </Container>
  );
}

export default App; 