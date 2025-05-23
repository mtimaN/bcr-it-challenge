
const adsByCluster = {
  0: [
    {
      title_ro: "Economisește mai ușor",
      title_eng: "Save easier",

      description_ro: "Vezi cele mai bune conturi de economii.",
      description_eng: "View the best savings accounts for you.",

      image: "images/ads/economyAccount.jpg"
    },
    {
      title_ro: "Termen lung, câștig mare",
      title_eng: "Consistent gain, long term",

      description_ro: "Depozite avantajoase pentru tine.",
      description_eng: "Profitable deposits for you.",

      image: "/images/ads/depositFavorable.jpg"
    },
    {
      title_ro: "Analizează-ți cheltuielile",
      title_eng: "Analyze your expenses",

      description_ro: "Grafice lunare pentru control financiar.",
      description_eng: "Monthly charts for financial control.",

      image: "/images/ads/financialGraphic.jpg"
    }
  ],
  1: [
    {
      title_ro: "Ai cheltuieli mari?",
      title_eng: "Are you a big spender?",

      description_ro: "Descoperă cardurile de credit smart.",
      description_eng: "Discover smart credit cards.",

      image: "images/ads/smartCard.jpg"
    },
    {
      title_ro: "Credit rapid, fără griji",
      title_eng: "Fast credit, worry-free",

      description_ro: "Aplică 100% online pentru împrumut.",
      description_eng: "Fully online application for loans.",

      image: "images/ads/loan.jpg "
    },
    {
      title_ro: "Folosește la maxim overdraft-ul",
      title_eng: "Maximize your overdraft benefits",

      description_ro: "Vezi limitele disponibile.",
      description_eng: "See available limits.",

      image: "images/ads/overdraft.jpg"
    }
  ],
  2: [
    {
      title_ro: "George te poate ajuta",
      title_eng: "George can help you",

      description_ro: "Descoperă beneficiile contului digital.",
      description_eng: "Discover the benefits of a digital account.",

      image: "images/ads/digitalAccount.jpg"
    },
    {
      title_ro: "Transferuri mai rapide",
      title_eng: "Faster transfers",

      description_ro: "Încearcă plățile instant.",
      description_eng: "Try instant payments.",

      image: "images/ads/contactless.jpg"
    },
    {
      title_ro: "Economisește fără efort",
      title_eng: "Save effortlessly",
  
      description_ro: "Setează un plan automat.",
      description_eng: "Set up a recurring plan.",

      image: "images/ads/financePlan.jpg"
    }
  ],
  3: [
    {
      title_ro: "Trimite bani rapid",
      title_eng: "Quick money transfer",

      description_ro: "Transfer instant către prieteni.",
      description_eng: "Instant transfer to friends.",

      image: "images/ads/fastTransfers.jpg"
    },
    {
      title_ro: "Economii simple",
      title_eng: "Simple savings",

      description_ro: "Economisește automat când cheltui.",
      description_eng: "Turn spending into saving automatically.",

      image: "images/ads/economyAuto.jpg"
    },
    {
      title_ro: "Cheltuie smart",
      title_eng: "Make smarter spending choices",

      description_ro: "Vezi unde se duc banii tăi.",
      description_eng: "See where your money goes.",
      
      image: "images/ads/moneyLocation.jpg"
    }
  ]
};

const allServices = { 
  transfers: { label_ro: "Transferuri", label_eng: "Transfers", color: "#d0f0fd", image: "images/services/transfersIcon.jpg" },
  deposits: { label_ro: "Depozite", label_eng: "Deposits", color: "#d1f7c4", image: "images/services/depositsIcon.jpg" },
  savings: { label_ro: "Economii", label_eng: "Savings", color: "#fceabb", image: "images/services/savingsIcon.jpg" },
  loans: { label_ro: "Împrumuturi", label_eng: "Loans", color: "#ffe0e9", image: "images/services/loanIcon.jpg" },
  cards: { label_ro: "Carduri", label_eng: "Cards", color: "#e1d7fc", image: "images/services/cardsIcon.jpg" },
  insights: { label_ro: "Analiză", label_eng: "Analysis", color: "#ffd8b1", image: "images/services/insightsIcon.jpg" },
};

const serviceOrderByCluster = {
  0: ['savings', 'deposits', 'insights', 'transfers', 'loans', 'cards'],
  1: ['cards', 'loans', 'transfers', 'deposits', 'insights', 'savings'],
  2: ['insights', 'transfers', 'savings', 'cards', 'loans', 'deposits'],
  3: ['transfers', 'savings', 'insights', 'loans', 'deposits', 'cards'],
};

export { adsByCluster, allServices, serviceOrderByCluster };
