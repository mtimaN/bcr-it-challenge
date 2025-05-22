import {
    LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid
  } from 'recharts';
  import './SpentChart.css';
  
  const SpentChart = ({ theme }) => {
    const isDark = theme === 'dark';
  
    const rawTransactions = [
        { day: 1, spending: 50 },
        { day: 2, spending: 0 },
        { day: 3, spending: 100 },
        { day: 4, spending: 20 },
        { day: 5, spending: 35 },
        { day: 6, spending: 500 },
        { day: 10, spending: 20 },
        { day: 12, spending: 200 },
        { day: 16, spending: 25 },
        { day: 21, spending: 30 },
        { day: 25, spending: 45 },
    ];

    const today = new Date();

    function getLastDayOfMonth(year, month) {
        return new Date(year, month + 1, 0).getDate();
    }


    function getCumulativeCheckpoints(data) {
        // Sort by day just in case
        const sortedData = [...data].sort((a, b) => a.day - b.day);
      
        let total = 0;
        return sortedData.map(entry => {
          total += entry.spending;
          return {
            day: entry.day,
            total
          };
        });
    }

    const formatBalance = (amount) => {
        return amount
          .toFixed(2)
          .toString()
          .replace('.', ',')
          .replace(/\B(?=(\d{3})+(?!\d))/g, '.');
        };
      
    const chartData = getCumulativeCheckpoints(rawTransactions);
    const checkpoints = [1, 8, 15, getLastDayOfMonth(today.getFullYear(), today.getMonth())];
  
    return (
      <>
        <div className={`chart-container ${isDark ? 'dark' : ''}`}>
            <h3 className={`chart-title ${isDark ? 'dark' : ''}`}>
                Luna asta
            </h3>
          <div className="chart-inner-wrapper">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={chartData} className="line-chart" margin={{ top: 50, right: 30, left: 63, bottom: 0 }}>
                <CartesianGrid className="grid" strokeDasharray="2 3" />

                <XAxis
                    dataKey="day"
                    type="number"
                    domain={['dataMin', 'dataMax']}
                    ticks={checkpoints}
                    tick={({ x, y, payload }) => (
                        <text
                        x={x}
                        y={y + 20}
                        textAnchor="middle"
                        className={`axis-tick ${isDark ? 'dark' : ''}`}
                        >
                        {payload.value}
                        </text>
                    )}
                />

                <YAxis
                    tick={({ x, y, payload }) => (
                        <text
                        x={x - 10}
                        y={y + 4}
                        textAnchor="end"
                        className={`axis-tick ${isDark ? 'dark' : ''}`}
                        >
                        {formatBalance(payload.value)}
                        </text>
                    )}
                />

                <Tooltip
                    wrapperClassName={`tooltip-wrapper ${isDark ? 'dark' : ''}`}
                    contentClassName="tooltip-content"
                    labelClassName="tooltip-label"
                    itemClassName="tooltip-item"
                    formatter={(value) => [`${formatBalance(value)} RON`, 'Cheltuit']}
                />

                <Line
                    type="linear"               // spiky line
                    dataKey="total"          // plot raw spending
                    className="spending-line"
                    strokeWidth={2}
                    dot
                    activeDot
                />
            </LineChart>
          </ResponsiveContainer>
          </div>
        </div>
      </>
    );
  };
  
  export default SpentChart;
  