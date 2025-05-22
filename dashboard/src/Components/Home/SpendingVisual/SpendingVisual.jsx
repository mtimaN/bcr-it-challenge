import React from 'react';
import './SpendingVisual.css';
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer
} from 'recharts';

const data = [
  { month: 'Ian 2025', spending: 400, earning: 200 },
  { month: 'Feb 2025', spending: 300, earning: 600 },
  { month: 'Mar 2025', spending: 500, earning: 100 },
  { month: 'Apr 2025', spending: 700, earning: 250 },
  { month: 'Mai 2025', spending: 600, earning: 410 },
  { month: 'Iun 2025', spending: 450, earning: 230 },
];

const SpendingVisual = ({ theme = 'light' }) => {
  const isDark = theme === 'dark';

  const formatBalance = (amount) => {
    return amount
      .toFixed(2)
      .toString()
      .replace('.', ',')
      .replace(/\B(?=(\d{3})+(?!\d))/g, '.');
    };

  const customTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      return (
        <div className={`tooltip-box ${isDark ? 'dark' : ''}`}>
          <p className={`tooltip-label ${isDark ? 'dark' : ''}`}>{label}</p>
          {payload.map((entry, index) => (
            <p
              key={index}
              className={`tooltip-entry ${isDark ? 'dark' : ''}`}
              style={{ '--entry-color': entry.color }}
            >
              {`${entry.name}: ${formatBalance(entry.value)} RON`}
            </p>
          ))}
        </div>
      );
    }
    return null;
  };

  const renderXAxisTick = ({ x, y, payload }) => (
    <text
      x={x}
      y={y + 20}  // Only vertical shift
      textAnchor="middle"
      className={`axis-tick ${isDark ? 'dark' : ''}`}
    >
      {payload.value}
    </text>
  );

  const renderYAxisTick = ({ x, y, payload }) => (
    <text
      x={x - 10}  // Only horizontal shift
      y={y + 4}
      textAnchor="end"
      className={`axis-tick ${isDark ? 'dark' : ''}`}
    >
      {`${formatBalance(payload.value)}`}
    </text>
  );

  return (
    <div className={`spending-container ${isDark ? 'dark' : ''}`}>
      <h3 className={`spending-title ${isDark ? 'dark' : ''}`}>
        Cheltuieli - Venituri Lunare
      </h3>

      <ResponsiveContainer width="100%" height="100%">
        <BarChart
          data={data}
          margin={{ top: 50, right: 30, left: 63, bottom: 0 }}
          barCategoryGap="20%"
          barGap={5}
        >
          <CartesianGrid
            strokeDasharray="3 3"
            stroke={isDark ? '#444' : '#e0e0e0'}
            horizontal
            vertical={false}
          />

        <XAxis
        dataKey="month"
        axisLine={{ stroke: isDark ? '#666' : '#333' }}
        tickLine={false}
        tick={renderXAxisTick}
        />

        <YAxis
        axisLine={{ stroke: isDark ? '#666' : '#333' }}
        tickLine={false}
        tick={renderYAxisTick}
        />

        <Tooltip content={customTooltip} />
        <Bar
        dataKey="spending"
        name="Cheltuit"
        fill="#b50606"
        radius={[4, 4, 0, 0]}
        animationDuration={800}
        animationEasing="ease-out"
        />
        <Bar
        dataKey="earning"
        name="Primit"
        fill="#186702"
        radius={[4, 4, 0, 0]}
        animationDuration={800}
        animationEasing="ease-out"
        />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};

export default SpendingVisual;
