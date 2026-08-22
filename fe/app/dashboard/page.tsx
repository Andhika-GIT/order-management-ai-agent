import { SectionCards } from "@/components/section-cards";

import { getDasboardData } from "../action/dashboard";

export default async function Page() {

  const data = await getDasboardData()

  return (
    <>
      <SectionCards total_products={data.total_products}/>
      {/* <div className="px-4 lg:px-6">
        <ChartAreaInteractive />
      </div> */}
    </>
  );
}
