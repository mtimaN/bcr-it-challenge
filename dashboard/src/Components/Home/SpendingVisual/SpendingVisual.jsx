import React from 'react';
import './SpendingVisual.css';
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer
} from 'recharts';

const data = [
  { month_ro: 'Dec 2024', month_eng: 'Dec 2024', spending: 400, earning: 200 },
  { month_ro: 'Ian 2025', month_eng: 'Jan 2025', spending: 300, earning: 600 },
  { month_ro: 'Feb 2025', month_eng: 'Feb 2025', spending: 500, earning: 100 },
  { month_ro: 'Mar 2025', month_eng: 'Mar 2025', spending: 700, earning: 250 },
  { month_ro: 'Apr 2025', month_eng: 'Apr 2025', spending: 600, earning: 410 },
  { month_ro: 'Mai 2025', month_eng: 'May 2025', spending: 450, earning: 230 },
];

const SpendingVisual = ({ theme, lang }) => {

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
        <div className={`tooltip-box ${theme === 'dark' ? 'dark' : ''}`}>
          <p className={`tooltip-label ${theme === 'dark' ? 'dark' : ''}`}>{label}</p>
          {payload.map((entry, index) => (
            <p
              key={index}
              className={`tooltip-entry ${theme === 'dark' ? 'dark' : ''}`}
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
      y={y + 20}
      textAnchor="middle"
      className={`axis-tick ${theme === 'dark' ? 'dark' : ''}`}
    >
      {payload.value}
    </text>
  );

  const renderYAxisTick = ({ x, y, payload }) => (
    <text
      x={x - 10}
      y={y + 4}
      textAnchor="end"
      className={`axis-tick ${theme === 'dark' ? 'dark' : ''}`}
    >
      {`${formatBalance(payload.value)}`}
    </text>
  );

  return (
    <div className={`spending-container ${theme === 'dark' ? 'dark' : ''}`}>
      <h3 className={`spending-title ${theme === 'dark' ? 'dark' : ''}`}>
        {lang === 'RO' ? 'Flux de numerar': 'Cash Flow'}
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
            stroke={theme === 'dark' ? '#444' : '#e0e0e0'}
            horizontal
            vertical={false}
          />

        <XAxis
        dataKey={lang === 'RO' ? 'month_ro' : 'month_eng'}
        axisLine={{ stroke: theme === 'dark' ? '#666' : '#333' }}
        tickLine={false}
        tick={renderXAxisTick}
        />

        <YAxis
        axisLine={{ stroke: theme === 'dark' ? '#666' : '#333' }}
        tickLine={false}
        tick={renderYAxisTick}
        />

        <Tooltip content={customTooltip} />
        <Bar
          dataKey="spending"
          name={lang === 'RO' ? 'Cheltuit' : 'Spent'}
          fill={theme === 'dark' ? '#ff7a7a' : '#b50606'}
          radius={[4, 4, 0, 0]}
          animationDuration={800}
          animationEasing="ease-out"
        />

        <Bar
          dataKey="earning"
          name={lang === 'RO' ? 'Primit' : 'Received'}
          fill={theme === 'dark' ? '#7aff7a' : '#186702'}
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
